package middleware

import (
	"fmt"

	"github.com/TicketsBot-cloud/common/premium"
	"github.com/Chaoskjell44/dashboard/config"
	"github.com/Chaoskjell44/dashboard/rpc"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
)

func VerifyWhitelabel(isApi bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId := ctx.Keys["userid"].(uint64)

		tier, err := rpc.PremiumClient.GetTierByUser(ctx, userId, false)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return
		}

		if tier < premium.Whitelabel {
			var isForced bool
			for _, id := range config.Conf.ForceWhitelabel {
				if id == userId {
					isForced = true
					break
				}
			}

			if !isForced {
				if isApi {
					ctx.AbortWithStatusJSON(402, gin.H{
						"success": false,
						"error":   "You must have the whitelabel premium tier",
					})
				} else {
					ctx.Redirect(302, fmt.Sprintf("%s/premium", config.Conf.Server.MainSite))
					ctx.Abort()
				}
				return
			}
		}
	}
}
