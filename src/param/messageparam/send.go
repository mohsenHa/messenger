package messageparam

import (
	"context"
	"time"
)

type SendRequest struct {
	Ctx     context.Context
	ToId    string `json:"to_id"`
	FromId  string
	Message string `json:"message"`
}

type SendResponse struct {
	SendMessage SendMessage `json:"send_message"`
}

type SendMessage struct {
	Id       string    `json:"id"`
	From     SendUser  `json:"from"`
	To       SendUser  `json:"to"`
	Body     string    `json:"body"`
	SendTime time.Time `json:"send_time"`
}

type SendUser struct {
	Id        string `json:"id"`
	PublicKey string `json:"public_key"`
}
