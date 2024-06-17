package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mohsenHa/messenger/logger"
	"io"
	"net/http"
	"time"
)

type GetPublicKeyRequest struct {
	ID    string `json:"id"`
	Token string `json:"hidden"`
}
type GetPublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

func GetPublicKey(req GetPublicKeyRequest) (GetPublicKeyResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return GetPublicKeyResponse{}, err
	}
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, targetHost.path("user/public_key"), bytes.NewBuffer(b))
	if err != nil {
		return GetPublicKeyResponse{}, err
	}
	request.Header.Set("Authorization", "Bearer "+req.Token)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return GetPublicKeyResponse{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetPublicKeyResponse{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		logger.NewLog("error on closing response body")
	}
	if resp.StatusCode != http.StatusOK {
		return GetPublicKeyResponse{}, fmt.Errorf("error: %v", string(body))
	}

	response := GetPublicKeyResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return GetPublicKeyResponse{}, err
	}

	return response, nil

}
