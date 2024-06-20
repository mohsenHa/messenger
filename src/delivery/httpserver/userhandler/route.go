package userhandler

import (
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/delivery/httpserver/middleware"
)

func (h Handler) SetRoutes(messageGroup *echo.Group) {
	messageGroup.POST("/register", h.userRegister)
	messageGroup.POST("/verify", h.userVerify)
	messageGroup.POST("/login", h.userLogin)
	messageGroup.POST("/id", h.userID)
	messageGroup.POST("/public_key", h.publicKey, middleware.Auth(h.authService))
	messageGroup.GET("/info", h.info, middleware.Auth(h.authService))
}
