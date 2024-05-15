package userservice

import (
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/param/user"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
)

func (s Service) Register(req user.RegisterRequest) (user.RegisterResponse, error) {
	activeCode := random.String(s.config.KeyLength)

	encryptedCode, err := encryptdecrypt.Encrypt(req.PublicKey, []byte(activeCode))
	if err != nil {
		return user.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	u := entity.User{
		PublicKey:  req.PublicKey,
		ActiveCode: activeCode,
		Status:     0,
	}

	// create new user in storage
	_, err = s.repo.Register(u)
	if err != nil {
		return user.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	// return created user
	return user.RegisterResponse{EncryptedCode: string(encryptedCode)}, nil
}
