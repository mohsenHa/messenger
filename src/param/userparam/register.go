package userparam

import "context"

type RegisterRequest struct {
	Ctx       context.Context
	PublicKey string `json:"public_key"`
}

type RegisterResponse struct {
	Id            string `json:"id"`
	EncryptedCode string `json:"encrypted_code"`
}
