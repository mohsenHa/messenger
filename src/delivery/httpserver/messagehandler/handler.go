package messagehandler

import (
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/validator/messagevalidator"
)

type Handler struct {
	messageSvc       messageservice.Service
	messageValidator messagevalidator.Validator
	authService      authservice.Service
}

func New(messageSvc messageservice.Service, authService authservice.Service, messageValidator messagevalidator.Validator) Handler {
	return Handler{
		messageSvc:       messageSvc,
		messageValidator: messageValidator,
		authService:      authService,
	}
}
