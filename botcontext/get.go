package botcontext

import (
	"context"
	"fmt"

	"github.com/TicketsBot-cloud/common/restcache"
	"github.com/Chaoskjell44/dashboard/config"
	dbclient "github.com/Chaoskjell44/dashboard/database"
	"github.com/Chaoskjell44/dashboard/redis"
	"github.com/rxdn/gdl/rest/ratelimit"
)

func ContextForGuild(guildId uint64) (*BotContext, error) {
	whitelabelBotId, isWhitelabel, err := dbclient.Client.WhitelabelGuilds.GetBotByGuild(context.Background(), guildId)
	if err != nil {
		return nil, err
	}

	if isWhitelabel {
		res, err := dbclient.Client.Whitelabel.GetByBotId(context.Background(), whitelabelBotId)
		if err != nil {
			return nil, err
		}

		rateLimiter := ratelimit.NewRateLimiter(ratelimit.NewRedisStore(redis.Client.Client, fmt.Sprintf("ratelimiter:%d", whitelabelBotId)), 1)

		return &BotContext{
			BotId:       res.BotId,
			Token:       res.Token,
			RateLimiter: rateLimiter,
			RestCache:   restcache.NewRedisRestCache(redis.Client.Client, res.Token, rateLimiter),
		}, nil
	} else {
		return PublicContext(), nil
	}
}

func PublicContext() *BotContext {
	rateLimiter := ratelimit.NewRateLimiter(ratelimit.NewRedisStore(redis.Client.Client, "ratelimiter:public"), 1)

	return &BotContext{
		BotId:       config.Conf.Bot.Id,
		Token:       config.Conf.Bot.Token,
		RateLimiter: rateLimiter,
		RestCache:   restcache.NewRedisRestCache(redis.Client.Client, config.Conf.Bot.Token, rateLimiter),
	}
}
