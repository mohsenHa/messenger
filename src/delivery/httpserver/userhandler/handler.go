package userhandler

import (
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/uservalidator"
)

type Handler struct {
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authService   authservice.Service
}

func New(userSvc userservice.Service,
	userValidator uservalidator.Validator,
	authService authservice.Service) Handler {
	return Handler{
		userSvc:       userSvc,
		userValidator: userValidator,
		authService:   authService,
	}
}
