package forms

import (
	"net/http"

	"github.com/Chaoskjell44/dashboard/app"
	dbclient "github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/TicketsBot-cloud/database"
	"github.com/gin-gonic/gin"
)

type createFormBody struct {
	Title string `json:"title"`
}

func CreateForm(c *gin.Context) {
	guildId := c.Keys["guildid"].(uint64)

	var data createFormBody
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, utils.ErrorStr("Invalid request body"))
		return
	}

	if len(data.Title) > 45 {
		c.JSON(400, utils.ErrorStr("Title is too long"))
		return
	}

	// 26^50 chance of collision
	customId, err := utils.RandString(30)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	id, err := dbclient.Client.Forms.Create(c, guildId, data.Title, customId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, app.NewServerError(err))
		return
	}

	form := database.Form{
		Id:       id,
		GuildId:  guildId,
		Title:    data.Title,
		CustomId: customId,
	}

	c.JSON(200, form)
}
