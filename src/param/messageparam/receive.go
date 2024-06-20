package messageparam

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReceiveRequest struct {
	Ctx         context.Context
	UserID      string
	EchoContext *echo.Context
	Request     *http.Request
	Response    *echo.Response
}

type ReceiveResponse struct {
	SendMessage SendMessage `json:"send_message"`
}
