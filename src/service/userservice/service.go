package userservice

import (
	"context"
	"github.com/mohsenHa/messenger/entity"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUserByPublicKey(ctx context.Context, publicKey string) (entity.User, error)
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
