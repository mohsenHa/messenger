package userhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/httpmsg"
	"github.com/mohsenHa/messenger/service/authservice"
)

func (h Handler) info(c echo.Context) error {
	var req userparam.InfoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if fieldErrors, err := h.userValidator.ValidateInfoRequest(req); err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	req.Ctx = c.Request().Context()

	u, ok := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid auth token",
		})
	}
	req.UserID = u.ID

	resp, err := h.userSvc.Info(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
