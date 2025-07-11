package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/TicketsBot-cloud/common/premium"
	"github.com/Chaoskjell44/dashboard/botcontext"
	dbclient "github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/rpc"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/rest"
	"github.com/rxdn/gdl/rest/request"
)

func ResendPanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	panelId, err := strconv.Atoi(ctx.Param("panelid"))
	if err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	// get existing
	panel, err := dbclient.Client.Panel.GetById(ctx, panelId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if panel.PanelId == 0 {
		ctx.JSON(404, utils.ErrorStr("Panel not found"))
		return
	}

	// check guild ID matches
	if panel.GuildId != guildId {
		ctx.JSON(403, utils.ErrorStr("Guild ID doesn't match"))
		return
	}

	if panel.ForceDisabled {
		ctx.JSON(400, utils.ErrorStr("This panel is disabled and cannot be modified: please reactivate premium to re-enable it"))
		return
	}

	// delete old message
	// TODO: Use proper context
	if err := rest.DeleteMessage(context.Background(), botContext.Token, botContext.RateLimiter, panel.ChannelId, panel.GuildId); err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && !unwrapped.IsClientError() {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}
	}

	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(ctx, guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	messageData := panelIntoMessageData(panel, premiumTier > premium.None)
	msgId, err := messageData.send(botContext)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			ctx.JSON(500, utils.ErrorStr("I do not have permission to send messages in the provided channel"))
		} else {
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	if err = dbclient.Client.Panel.UpdateMessageId(ctx, panel.PanelId, msgId); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, utils.SuccessResponse)
}
