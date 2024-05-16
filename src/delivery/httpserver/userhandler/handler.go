package userhandler

import (
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/uservalidator"
)

type Handler struct {
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(userSvc userservice.Service,
	userValidator uservalidator.Validator) Handler {
	return Handler{
		userSvc:       userSvc,
		userValidator: userValidator,
	}
}
