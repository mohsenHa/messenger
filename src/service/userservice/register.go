package userservice

import (
	"fmt"

	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) Register(req userparam.RegisterRequest) (userparam.RegisterResponse, error) {
	unexpectedError := "unexpected error: %w"

	u := entity.User{
		ID:        s.keyGen.CreateUserID(req.PublicKey),
		PublicKey: req.PublicKey,
		Code:      "",
		Status:    0,
	}
	// create new user in storage
	_, err := s.repo.Register(req.Ctx, u)
	if err != nil {
		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}
	encryptedCode, err := s.updateCode(req.Ctx, u)
	if err != nil {
		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}

	// return created user
	return userparam.RegisterResponse{
		EncryptedCode: encryptedCode,
		ID:            u.ID,
	}, nil
}
