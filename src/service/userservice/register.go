package userservice

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
)

func (s Service) Register(req userparam.RegisterRequest) (userparam.RegisterResponse, error) {
	activeCode := random.String(s.config.KeyLength)
	publicKey, err := base64.RawStdEncoding.DecodeString(req.PublicKey)
	unexpectedError := "unexpected error: %w"
	if err != nil {
		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}
	encryptedCode, err := encryptdecrypt.Encrypt(publicKey, []byte(activeCode))
	if err != nil {
		fmt.Println(err)

		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}

	u := entity.User{
		PublicKey:  req.PublicKey,
		ActiveCode: activeCode,
		Status:     0,
	}
	// create new user in storage
	_, err = s.repo.Register(u)
	if err != nil {
		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}
	// return created user
	return userparam.RegisterResponse{
			EncryptedCode: base64.RawStdEncoding.EncodeToString(encryptedCode)},
		nil
}
