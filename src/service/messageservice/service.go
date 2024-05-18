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
	GetUserById(ctx context.Context, id string) (entity.User, error)
}

func New(rabbitmq *rabbitmq.ChannelAdapter, userRepo UserRepo) Service {
	return Service{
		rabbitmq: rabbitmq,
		userRepo: userRepo,
	}
}
