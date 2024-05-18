package entity

import "time"

type Message struct {
	Id       string
	From     User
	To       User
	Body     string
	SendTime time.Time
}
