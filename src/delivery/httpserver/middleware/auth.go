package middleware

import (
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/service/authservice"
)

func Auth(service authservice.Service) echo.MiddlewareFunc {
	publicKey, err := service.GetPublicKey()
	if err != nil {
		panic(err)
	}

	return mw.WithConfig(mw.Config{
		ContextKey:    config.AuthMiddlewareContextKey,
		SigningKey:    publicKey,
		SigningMethod: service.GetSigningMethod().Name,
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}
			return claims, nil
		},
	})
}
