package userparam

import "context"

type IDRequest struct {
	Ctx       context.Context
	PublicKey string `json:"public_key"`
}

type IDResponse struct {
	ID string `json:"id"`
}
