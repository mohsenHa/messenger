package messageparam

import (
	"context"
	"time"
)

type SendRequest struct {
	Ctx     context.Context
	ToID    string `json:"to_id"`
	FromID  string
	Message string `json:"message"`
}

type SendResponse struct {
	SendMessage SendMessage `json:"send_message"`
}

type SendMessageType string

const (
	SendMessageDeliverType SendMessageType = "deliver"
	SendMessageMessageType SendMessageType = "message"
)

type SendMessage struct {
	ID       string          `json:"id"`
	Type     SendMessageType `json:"type"`
	From     SendUser        `json:"from"`
	To       SendUser        `json:"to"`
	Body     string          `json:"body"`
	SendTime time.Time       `json:"send_time"`
}

type SendUser struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
}
