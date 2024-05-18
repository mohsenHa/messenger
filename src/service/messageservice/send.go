package messageservice

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mohsenHa/messenger/param/messageparam"
	"time"
)

func (s Service) Send(req messageparam.SendRequest) (messageparam.SendResponse, error) {

	fromUser, _ := s.userRepo.GetUserById(req.Ctx, req.FromId)
	toUser, _ := s.userRepo.GetUserById(req.Ctx, req.ToId)

	id := uuid.New()
	sendMessage := messageparam.SendMessage{
		Id: id.String(),
		From: messageparam.SendUser{
			Id:        fromUser.Id,
			PublicKey: fromUser.PublicKey,
		},
		To: messageparam.SendUser{
			Id:        toUser.Id,
			PublicKey: toUser.PublicKey,
		},
		Body:     req.Message,
		SendTime: time.Now().UTC(),
	}
	messageByte, err := json.Marshal(sendMessage)
	if err != nil {
		return messageparam.SendResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	chanel := s.rabbitmq.GetInputChannel(req.ToId)
	chanel <- messageByte

	return messageparam.SendResponse{SendMessage: sendMessage}, nil
}
