package userservice

import (
	"context"

	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
)

func (s Service) updateCode(ctx context.Context, u entity.User) (string, error) {
	code, err := s.keyGen.CreateCode()
	if err != nil {
		return "", err
	}
	err = s.repo.UpdateCode(ctx, u.ID, encryptdecrypt.GetMD5Hash(code))
	if err != nil {
		return "", err
	}
	encryptedCode, err := s.keyGen.EncryptCode(code, u.PublicKey)
	if err != nil {
		return "", err
	}

	return encryptedCode, nil
}
