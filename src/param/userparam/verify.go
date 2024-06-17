package userparam

import "context"

type VerifyRequest struct {
	Ctx  context.Context
	ID   string `json:"id"`
	Code string `json:"code"`
}

type VerifyResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}
