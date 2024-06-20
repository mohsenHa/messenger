package userhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/httpmsg"
)

func (h Handler) userID(c echo.Context) error {
	var req userparam.IDRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if fieldErrors, err := h.userValidator.ValidateIDRequest(req); err != nil {
		msg, code := httpmsg.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}
	req.Ctx = c.Request().Context()
	resp, err := h.userSvc.ID(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
