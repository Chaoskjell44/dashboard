package api

import (
	"context"
	"strconv"
	"strings"

	dbclient "github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/rpc/cache"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
)

const pageLimit = 20
const builtInCount = 1

type (
	integrationWithMetadata struct {
		integrationResponse
		Author     *integrationAuthor `json:"author"`
		GuildCount int                `json:"guild_count"`
		Added      bool               `json:"added"`
	}

	integrationAuthor struct {
		Id            uint64             `json:"id,string"`
		Username      string             `json:"username"`
		Discriminator user.Discriminator `json:"discriminator"`
		Avatar        user.Avatar        `json:"avatar"`
	}
)

func ListIntegrationsHandler(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)
	guildId := ctx.Keys["guildid"].(uint64)

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 1 {
		page = 1
	}

	page -= 1

	limit := pageLimit
	if page == 0 {
		limit -= builtInCount
	}

	availableIntegrations, err := dbclient.Client.CustomIntegrationGuilds.GetAvailableIntegrationsWithActive(ctx, guildId, userId, limit, page*pageLimit)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	var authorIds []uint64
	integrations := make([]integrationWithMetadata, len(availableIntegrations))
	for i, integration := range availableIntegrations {
		var proxyToken *string
		if integration.ImageUrl != nil {
			tmp, err := utils.GenerateImageProxyToken(*integration.ImageUrl)
			if err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}

			proxyToken = &tmp
		}

		integrations[i] = integrationWithMetadata{
			integrationResponse: integrationResponse{
				Id:               integration.Id,
				OwnerId:          integration.OwnerId,
				WebhookHost:      utils.SecondLevelDomain(utils.GetUrlHost(strings.ReplaceAll(integration.WebhookUrl, "%", ""))),
				Name:             integration.Name,
				Description:      integration.Description,
				ImageUrl:         integration.ImageUrl,
				ProxyToken:       proxyToken,
				PrivacyPolicyUrl: integration.PrivacyPolicyUrl,
				Public:           integration.Public,
				Approved:         integration.Approved,
			},
			GuildCount: integration.GuildCount,
			Added:      integration.Active,
		}

		authorIds = append(authorIds, integration.OwnerId)
	}

	// Get author data for the integrations
	// TODO: Use proper context
	authors, err := cache.Instance.GetUsers(context.Background(), authorIds)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	for i, integration := range integrations {
		author, ok := authors[integration.OwnerId]
		if ok {
			integrations[i].Author = &integrationAuthor{
				Id:       author.Id,
				Username: author.Username,
				Avatar:   author.Avatar,
			}
		}
	}

	// Don't serve null
	if integrations == nil {
		integrations = make([]integrationWithMetadata, 0)
	}

	ctx.JSON(200, integrations)
}
