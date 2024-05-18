package messagehandler

import (
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/param/messageparam"
	"github.com/mohsenHa/messenger/pkg/httpmsg"
	"github.com/mohsenHa/messenger/service/authservice"
	"net/http"
)

func (h Handler) receiveMessage(c echo.Context) error {
	var req messageparam.ReceiveRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if fieldErrors, err := h.messageValidator.ValidateReceiveRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	req.Ctx = c.Request().Context()

	req.EchoContext = &c

	req.Response = c.Response()
	req.Request = c.Request()

	req.UserId = c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims).Id

	err := h.messageSvc.Receive(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
