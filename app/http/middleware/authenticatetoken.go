package middleware

import (
	"fmt"
	"strconv"

	"github.com/Chaoskjell44/dashboard/config"
	"github.com/Chaoskjell44/dashboard/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthenticateToken(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")

	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Conf.Server.Secret), nil
	})

	if err != nil {
		ctx.AbortWithStatusJSON(401, utils.ErrorJson(err))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, hasUserId := claims["userid"]
		if !hasUserId {
			ctx.AbortWithStatusJSON(401, utils.ErrorStr("Token is invalid"))
			return
		}

		parsedId, err := strconv.ParseUint(userId.(string), 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(401, utils.ErrorStr("Token is invalid"))
			return
		}

		if ctx.Keys == nil {
			ctx.Keys = make(map[string]interface{})
		}

		ctx.Keys["userid"] = parsedId
	} else {
		ctx.AbortWithStatusJSON(401, utils.ErrorStr("Token is invalid"))
		return
	}
}
