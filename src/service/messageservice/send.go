package messageservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mohsenHa/messenger/param/messageparam"
)

func (s Service) Send(req messageparam.SendRequest) (messageparam.SendResponse, error) {
	fromUser, err := s.userRepo.GetUserByID(req.Ctx, req.FromID)
	if err != nil {
		return messageparam.SendResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	toUser, err := s.userRepo.GetUserByID(req.Ctx, req.ToID)
	if err != nil {
		return messageparam.SendResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	id := uuid.New()
	sendMessage := messageparam.SendMessage{
		ID:   id.String(),
		Type: messageparam.SendMessageMessageType,
		From: messageparam.SendUser{
			ID:        fromUser.ID,
			PublicKey: fromUser.PublicKey,
		},
		To: messageparam.SendUser{
			ID:        toUser.ID,
			PublicKey: toUser.PublicKey,
		},
		Body:     req.Message,
		SendTime: time.Now().UTC(),
	}
	messageByte, err := json.Marshal(sendMessage)
	if err != nil {
		return messageparam.SendResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	chanel, err := s.rabbitmq.GetInputChannel(req.ToID)
	if err != nil {
		return messageparam.SendResponse{}, err
	}
	chanel <- messageByte

	return messageparam.SendResponse{SendMessage: sendMessage}, nil
}
