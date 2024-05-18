package userservice

import (
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) Id(req userparam.IdRequest) (userparam.IdResponse, error) {
	// return created user
	return userparam.IdResponse{
		Id: s.keyGen.CreateUserId(req.PublicKey),
	}, nil
}
