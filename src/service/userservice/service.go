package userservice

import (
	"context"
	"github.com/mohsenHa/messenger/entity"
)

type Repository interface {
	Register(ctx context.Context, u entity.User) (entity.User, error)
	Activate(ctx context.Context, id string) error
	GetUserById(ctx context.Context, id string) (entity.User, error)
}
type Config struct {
	KeyLength uint8 `koanf:"key_length"`
}

type Service struct {
	repo   Repository
	config Config
}

func New(repo Repository, config Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}
