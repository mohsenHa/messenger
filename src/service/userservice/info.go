package userservice

import (
	"github.com/mohsenHa/messenger/param/userparam"
)

func (s Service) Info(req userparam.InfoRequest) (userparam.InfoResponse, error) {
	user, err := s.repo.GetUserById(req.Ctx, req.UserId)
	if err != nil {
		return userparam.InfoResponse{}, err
	}

	// return created user
	return userparam.InfoResponse{
		Info: userparam.UserInfo{
			Id:        user.Id,
			Status:    user.Status,
			PublicKey: user.PublicKey,
		},
	}, nil
}
