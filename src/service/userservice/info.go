package userservice

import (
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) Info(req userparam.InfoRequest) (userparam.InfoResponse, error) {
	user, err := s.repo.GetUserByID(req.Ctx, req.UserID)
	if err != nil {
		return userparam.InfoResponse{}, err
	}

	// return created user
	return userparam.InfoResponse{
		Info: userparam.UserInfo{
			ID:        user.ID,
			Status:    user.Status,
			PublicKey: user.PublicKey,
		},
	}, nil
}
