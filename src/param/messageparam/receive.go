package messageparam

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ReceiveRequest struct {
	Ctx         context.Context
	UserId      string
	EchoContext *echo.Context
	Request     *http.Request
	Response    *echo.Response
}

type ReceiveResponse struct {
	SendMessage SendMessage `json:"send_message"`
}
