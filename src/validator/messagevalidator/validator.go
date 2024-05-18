package messagevalidator

import "github.com/mohsenHa/messenger/service/keygenerator"

type UserRepository interface {
	IsIdExist(publicKey string) (bool, error)
}

type Validator struct {
	urepo UserRepository
}

func New(repo UserRepository, keyGen keygenerator.Service) Validator {
	return Validator{urepo: repo}
}
