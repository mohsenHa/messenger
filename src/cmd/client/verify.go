package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VerifyRequest struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}
type VerifyResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

func Verify(request VerifyRequest) (VerifyResponse, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return VerifyResponse{}, err
	}

	resp, err := http.Post(targetHost.path("user/verify"), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return VerifyResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return VerifyResponse{}, fmt.Errorf("somthing failed: %+v", resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return VerifyResponse{}, err
	}
	rV := VerifyResponse{}
	err = json.Unmarshal(body, &rV)
	if err != nil {
		return VerifyResponse{}, err
	}
	return rV, nil
}
