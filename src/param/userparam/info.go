package userparam

import "context"

type UserInfo struct {
	Id        string `json:"id"`
	Status    uint   `json:"status"`
	PublicKey string `json:"public_key"`
}

type InfoRequest struct {
	Ctx    context.Context
	UserId string
}

type InfoResponse struct {
	Info UserInfo `json:"info"`
}
