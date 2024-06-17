package messagevalidator

import "github.com/mohsenHa/messenger/service/keygenerator"

type UserRepository interface {
	IsIDExist(publicKey string) (bool, error)
}

type Validator struct {
	urepo UserRepository
}

func New(repo UserRepository, _ keygenerator.Service) Validator {
	return Validator{urepo: repo}
}
