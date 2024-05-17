package userservice

import (
	"fmt"
	"github.com/mohsenHa/messenger/param/userparam"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
	"github.com/mohsenHa/messenger/pkg/errmsg"
)

func (s Service) Verify(req userparam.VerifyRequest) (userparam.VerifyResponse, error) {
	code := req.Code

	unexpectedError := "unexpected error: %w"
	u, err := s.repo.GetUserById(req.Ctx, req.Id)
	if err != nil {
		return userparam.VerifyResponse{}, fmt.Errorf(unexpectedError, err)
	}

	if u.Code != encryptdecrypt.GetMD5Hash(code) {
		return userparam.VerifyResponse{}, fmt.Errorf(errmsg.ErrorMsgInvalidCode)
	}

	err = s.repo.Activate(req.Ctx, req.Id)
	if err != nil {
		return userparam.VerifyResponse{}, fmt.Errorf(unexpectedError, err)
	}
	token, err := s.auth.CreateAccessToken(u)
	if err != nil {
		return userparam.VerifyResponse{}, fmt.Errorf(unexpectedError, err)
	}

	return userparam.VerifyResponse{
		Id:    u.Id,
		Token: token,
	}, nil
}
