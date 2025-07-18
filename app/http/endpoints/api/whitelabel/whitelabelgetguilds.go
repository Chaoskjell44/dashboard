package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Chaoskjell44/dashboard/app"
	"github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/rpc/cache"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
	cache2 "github.com/rxdn/gdl/cache"
)

func WhitelabelGetGuilds(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	bot, err := database.Client.Whitelabel.GetByUserId(c, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	// id -> name
	if bot.BotId == 0 {
		c.JSON(400, utils.ErrorStr("Whitelabel bot not found"))
		return
	}

	ids, err := database.Client.WhitelabelGuilds.GetGuilds(c, bot.BotId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	guilds := make(map[string]string)
	for i, id := range ids {
		if i >= 10 {
			idStr := strconv.FormatUint(id, 10)
			guilds[idStr] = idStr
			continue
		}

		// get guild name
		// TODO: Use proper context
		guild, err := cache.Instance.GetGuild(context.Background(), id)
		if err != nil {
			if errors.Is(err, cache2.ErrNotFound) {
				continue
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
				return
			}
		}

		guilds[strconv.FormatUint(id, 10)] = guild.Name
	}

	c.JSON(200, gin.H{
		"success": true,
		"guilds":  guilds,
	})
}
