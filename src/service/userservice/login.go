package userservice

import (
	"fmt"
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) Login(req userparam.LoginRequest) (userparam.LoginResponse, error) {
	u, err := s.repo.GetUserByID(req.Ctx, req.ID)
	fmt.Println(req.ID, u, err)
	unexpectedError := "unexpected error: %w"
	if err != nil {
		return userparam.LoginResponse{}, fmt.Errorf(unexpectedError, err)
	}
	encryptedCode, err := s.updateCode(req.Ctx, u)
	if err != nil {
		return userparam.LoginResponse{}, fmt.Errorf(unexpectedError, err)
	}
	// return created user
	return userparam.LoginResponse{
		EncryptedCode: encryptedCode,
	}, nil
}
