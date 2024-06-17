package userparam

import "context"

type PublicKeyRequest struct {
	Ctx context.Context
	ID  string `json:"id"`
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}
