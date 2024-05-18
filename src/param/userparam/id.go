package userparam

import "context"

type IdRequest struct {
	Ctx       context.Context
	PublicKey string `json:"public_key"`
}

type IdResponse struct {
	Id string `json:"id"`
}
