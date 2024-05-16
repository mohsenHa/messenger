package uservalidator

type Repository interface {
	IsPublicKeyUnique(publicKey string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}
