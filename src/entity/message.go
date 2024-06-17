package entity

import "time"

type Message struct {
	ID       string
	From     User
	To       User
	Body     string
	SendTime time.Time
}
