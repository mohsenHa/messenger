package userhandler

import (
	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/httpmsg"
	"net/http"
)

func (h Handler) userVerify(c echo.Context) error {
	var req userparam.VerifyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	req.Ctx = c.Request().Context()

	if fieldErrors, err := h.userValidator.ValidateVerifyRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.userSvc.Verify(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
