package messagehandler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/param/messageparam"
	"github.com/mohsenHa/messenger/pkg/httpmsg"
	"github.com/mohsenHa/messenger/service/authservice"
	"net/http"
)

func (h Handler) sendMessage(c echo.Context) error {
	var req messageparam.SendRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if fieldErrors, err := h.messageValidator.ValidateSendRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	req.Ctx = c.Request().Context()

	req.FromId = c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims).Id

	if req.FromId == req.ToId {
		msg, code := httpmsg.Error(fmt.Errorf("sender and receiver cannot be the same %s and %s",
			req.FromId, req.ToId))
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	resp, err := h.messageSvc.Send(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
