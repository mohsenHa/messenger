package userservice

import (
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) ID(req userparam.IDRequest) (userparam.IDResponse, error) {
	// return created user
	return userparam.IDResponse{
		ID: s.keyGen.CreateUserID(req.PublicKey),
	}, nil
}
