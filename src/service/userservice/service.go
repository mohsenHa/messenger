package userservice

import (
	"context"
	"github.com/mohsenHa/messenger/entity"
)

type Repository interface {
	Register(ctx context.Context, u entity.User) (entity.User, error)
	UpdateCode(ctx context.Context, id, code string) error
	Activate(ctx context.Context, id string) error
	GetUserByID(ctx context.Context, id string) (entity.User, error)
}
type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
}

type KeyGenerator interface {
	CreateCode() (string, error)
	EncryptCode(code, publicKey string) (string, error)
	CreateUserID(publicKey string) string
}

type Service struct {
	repo   Repository
	auth   AuthGenerator
	keyGen KeyGenerator
}

func New(repo Repository, auth AuthGenerator, keyGen KeyGenerator) Service {
	return Service{
		repo:   repo,
		auth:   auth,
		keyGen: keyGen,
	}
}
