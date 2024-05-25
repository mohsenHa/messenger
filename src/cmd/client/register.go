package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RegisterRequest struct {
	PublicKey string `json:"public_key"`
}
type RegisterResponse struct {
	Id                string `json:"id"`
	EncryptedCode     string `json:"encrypted_code"`
	EncryptedCodeByte []byte
}

func Register(request RegisterRequest) (RegisterResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return RegisterResponse{}, err
	}
	resp, err := http.Post(targetHost.path("user/register"), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return RegisterResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisterResponse{}, err
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
