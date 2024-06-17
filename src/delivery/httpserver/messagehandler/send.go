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

	u, ok := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid auth token",
		})
	}
	req.FromID = u.ID
	if req.FromID == req.ToID {
		msg, code := httpmsg.Error(fmt.Errorf("sender and receiver cannot be the same %s and %s",
			req.FromID, req.ToID))

		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	resp, err := h.messageSvc.Send(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
