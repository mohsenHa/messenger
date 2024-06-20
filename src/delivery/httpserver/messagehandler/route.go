package messagehandler

import (
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/delivery/httpserver/middleware"
)

func (h Handler) SetRoutes(messageGroup *echo.Group) {
	messageGroup.POST("/send", h.sendMessage, middleware.Auth(h.authService))
	messageGroup.GET("/receive", h.receiveMessage, middleware.Auth(h.authService))
}
