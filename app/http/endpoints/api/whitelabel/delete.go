package api

import (
	"net/http"

	"github.com/TicketsBot-cloud/common/whitelabeldelete"
	"github.com/Chaoskjell44/dashboard/app"
	"github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/redis"
	"github.com/gin-gonic/gin"
)

func WhitelabelDelete(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	// Check if this is a different token
	botId, err := database.Client.Whitelabel.Delete(c, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	if botId != nil {
		// TODO: Kafka
		go whitelabeldelete.Publish(redis.Client.Client, *botId)

	}

	c.Status(http.StatusNoContent)
}
