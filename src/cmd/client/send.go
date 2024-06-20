package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mohsenHa/messenger/param/messageparam"
)

type SendRequest struct {
	Message string `json:"message"`
	ToID    string `json:"to_id"`
	Token   string `json:"-"`
}
type SendResponse struct {
	SendMessage messageparam.SendMessage `json:"send_message"`
}

func Send(req SendRequest) (messageparam.SendMessage, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return messageparam.SendMessage{}, err
	}
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetHost.path("message/send"), bytes.NewBuffer(b))
	if err != nil {
		return messageparam.SendMessage{}, err
	}
	request.Header.Set("Authorization", "Bearer "+req.Token)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return messageparam.SendMessage{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return messageparam.SendMessage{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		return messageparam.SendMessage{}, fmt.Errorf("error: %v", string(body))
	}
	response := SendResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return messageparam.SendMessage{}, err
	}

	return response.SendMessage, nil
}
