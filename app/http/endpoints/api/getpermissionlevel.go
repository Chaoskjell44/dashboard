package api

import (
	"net/http"
	"strconv"

	"github.com/Chaoskjell44/dashboard/app"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func GetPermissionLevel(c *gin.Context) {
	userId := c.Keys["userid"].(uint64)

	guildId, err := strconv.ParseUint(c.Query("guild"), 10, 64)
	if err != nil {
		c.JSON(400, utils.ErrorStr("Invalid guild ID"))
		return
	}

	// TODO: Use proper context
	permissionLevel, err := utils.GetPermissionLevel(c, guildId, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	c.JSON(200, gin.H{
		"success":          true,
		"permission_level": permissionLevel,
	})
}
