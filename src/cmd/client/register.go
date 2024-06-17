package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RegisterRequest struct {
	PublicKey string `json:"public_key"`
}
type RegisterResponse struct {
	ID                string `json:"id"`
	EncryptedCode     string `json:"encrypted_code"`
	EncryptedCodeByte []byte
}

func Register(request RegisterRequest) (RegisterResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return RegisterResponse{}, err
	}
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetHost.path("user/register"), bytes.NewBuffer(b))
	if err != nil {
		return RegisterResponse{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return RegisterResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisterResponse{}, err
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode != http.StatusCreated {
		return RegisterResponse{}, fmt.Errorf("error: %v", string(body))
	}

	response := RegisterResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return RegisterResponse{}, err
	}
	response.EncryptedCodeByte, err = base64.RawStdEncoding.DecodeString(response.EncryptedCode)
	if err != nil {
		return RegisterResponse{}, err
	}

	return response, nil

}
