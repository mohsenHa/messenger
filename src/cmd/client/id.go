package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetPublicKeyRequest struct {
	Id    string `json:"id"`
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
	request, err := http.NewRequest(http.MethodPost, targetHost.path("user/public_key"), bytes.NewBuffer(b))
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetPublicKeyResponse{}, err
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
