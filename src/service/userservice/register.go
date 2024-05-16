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
	publicKey, err := base64.StdEncoding.DecodeString(req.PublicKey)
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
		Id:         encryptdecrypt.GetMD5Hash(req.PublicKey),
		PublicKey:  req.PublicKey,
		ActiveCode: encryptdecrypt.GetMD5Hash(activeCode),
		Status:     0,
	}
	// create new user in storage
	_, err = s.repo.Register(req.Ctx, u)
	if err != nil {
		return userparam.RegisterResponse{}, fmt.Errorf(unexpectedError, err)
	}
	// return created user
	return userparam.RegisterResponse{
		EncryptedCode: base64.RawStdEncoding.EncodeToString(encryptedCode),
		Id:            u.Id,
	}, nil
}
