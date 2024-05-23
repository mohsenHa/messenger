package userservice

import (
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) PublicKey(req userparam.PublicKeyRequest) (userparam.PublicKeyResponse, error) {
	user, err := s.repo.GetUserById(req.Ctx, req.Id)
	if err != nil {
		return userparam.PublicKeyResponse{}, err
	}

	// return created user
	return userparam.PublicKeyResponse{
		PublicKey: user.PublicKey,
	}, nil
}
