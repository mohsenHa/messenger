package userhandler

import "github.com/labstack/echo/v4"

func (h Handler) SetRoutes(messageGroup *echo.Group) {

	messageGroup.POST("/register", h.userRegister)
	messageGroup.POST("/verify", h.userVerify)
	messageGroup.POST("/login", h.userLogin)
	messageGroup.POST("/id", h.userId)

}
