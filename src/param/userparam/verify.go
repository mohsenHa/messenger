package userparam

import "context"

type VerifyRequest struct {
	Ctx  context.Context
	Id   string `json:"id"`
	Code string `json:"code"`
}

type VerifyResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
