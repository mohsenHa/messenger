package uservalidator

import "github.com/mohsenHa/messenger/service/keygenerator"

type Repository interface {
	IsIDUnique(id string) (bool, error)
	IsIDExist(id string) (bool, error)
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
