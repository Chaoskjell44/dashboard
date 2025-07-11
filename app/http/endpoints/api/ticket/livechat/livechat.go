package livechat

import (
	"net/http"
	"strconv"

	"github.com/Chaoskjell44/dashboard/config"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == config.Conf.Server.BaseUrl
	},
}

func GetLiveChatHandler(sm *SocketManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		guildId, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, utils.ErrorJson(err))
			return
		}

		ticketId, err := strconv.Atoi(c.Param("ticketId"))
		if err != nil {
			c.JSON(400, utils.ErrorJson(err))
			return
		}

		client := NewClient(sm, conn, c, guildId, ticketId)
		sm.register <- client
		go client.StartReadLoop()
		go client.StartWriteLoop()
	}
}
