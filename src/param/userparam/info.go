package userparam

import "context"

type UserInfo struct {
	ID        string `json:"id"`
	Status    uint   `json:"status"`
	PublicKey string `json:"public_key"`
}

type InfoRequest struct {
	Ctx    context.Context
	UserID string
}

type InfoResponse struct {
	Info UserInfo `json:"info"`
}
