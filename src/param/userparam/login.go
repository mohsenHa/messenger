package userparam

import "context"

type LoginRequest struct {
	Ctx context.Context
	ID  string `json:"id"`
}

type LoginResponse struct {
	EncryptedCode string `json:"encrypted_code"`
}
