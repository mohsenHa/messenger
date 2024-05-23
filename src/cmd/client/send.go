package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mohsenHa/messenger/param/messageparam"
	"io"
	"net/http"
)

type SendRequest struct {
	Message string `json:"message"`
	ToId    string `json:"to_id"`
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
	request, err := http.NewRequest(http.MethodPost, targetHost.path("message/send"), bytes.NewBuffer(b))
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return messageparam.SendMessage{}, err
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
