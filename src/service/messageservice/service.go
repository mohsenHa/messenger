package messageservice

import (
	"context"

	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/entity"
)

type Service struct {
	rabbitmq *rabbitmq.ChannelAdapter
	userRepo UserRepo
}

type UserRepo interface {
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}

func New(rmq *rabbitmq.ChannelAdapter, userRepo UserRepo) Service {
	return Service{
		rabbitmq: rmq,
		userRepo: userRepo,
	}
}
