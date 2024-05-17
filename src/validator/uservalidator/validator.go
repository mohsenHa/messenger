package uservalidator

import "github.com/mohsenHa/messenger/service/keygenerator"

type Repository interface {
	IsIdUnique(publicKey string) (bool, error)
}

type Validator struct {
	repo   Repository
	keyGen keygenerator.Service
}

func New(repo Repository, keyGen keygenerator.Service) Validator {
	return Validator{
		repo:   repo,
		keyGen: keyGen,
	}
}
