package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LoginRequest struct {
	ID string `json:"id"`
}
type LoginResponse struct {
	EncryptedCode     string `json:"encrypted_code"`
	EncryptedCodeByte []byte
}

type VerifyRequest struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}
type VerifyResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type User struct {
	ID            string          `json:"id"`
	Token         string          `json:"token"`
	PublicKey     string          `json:"public_key"`
	PrivateKey    string          `json:"private_key"`
	PrivateKeyRSA *rsa.PrivateKey `json:"-"`
	UserFile      string          `json:"-"`
}

const ContentTypeApplicationJSON = "application/json"

func New(userFile string) (User, error) {
	keySize := 512
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return User{}, err
	}

	err = privateKey.Validate()
	if err != nil {
		return User{}, err
	}

	// Extract public component.
	pub := privateKey.Public()

	// Encode public key to PKCS#1 ASN.1 PEM.
	byt, err := x509.MarshalPKIXPublicKey(pub.(*rsa.PublicKey))
	if err != nil {
		return User{}, err
	}
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: byt,
		},
	)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(pubPEM)

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(keyPEM)

	rR, err := Register(RegisterRequest{PublicKey: publicKeyBase64})
	if err != nil {
		return User{}, err
	}
	user := User{
		ID:            rR.ID,
		Token:         "",
		PublicKey:     publicKeyBase64,
		PrivateKey:    privateKeyBase64,
		PrivateKeyRSA: privateKey,
		UserFile:      userFile,
	}

	decryptedBytes, err := user.Decrypt(rR.EncryptedCodeByte)
	if err != nil {
		return User{}, err
	}
	err = user.verify(string(decryptedBytes))
	if err != nil {
		return User{}, err
	}
	user.Store()

	return user, nil
}

func NewUserFromFile(userFile string) (User, error) {
	user := User{
		UserFile: userFile,
	}
	file, err := os.Open(userFile)
	if err != nil {
		return New(userFile)
	}

	j, err := io.ReadAll(file)
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(j, &user)
	if err != nil {
		return User{}, err
	}
	keyPEMbyte, err := base64.StdEncoding.DecodeString(user.PrivateKey)
	if err != nil {
		return User{}, err
	}
	keyPEM, _ := pem.Decode(keyPEMbyte)
	key, err := x509.ParsePKCS1PrivateKey(keyPEM.Bytes)
	if err != nil {
		return User{}, err
	}
	user.PrivateKeyRSA = key

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	check, err := user.Check()
	if err != nil {
		return User{}, err
	}
	if !check {
		err := user.Login()
		if err != nil {
			return User{}, err
		}
	}

	return user, nil
}

func (u *User) GetPrivateKey() (*rsa.PrivateKey, error) {
	keyPEM, err := base64.StdEncoding.DecodeString(u.PrivateKey)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyPEM)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (u *User) Decrypt(encryptedByte []byte) ([]byte, error) {
	privateKey, err := u.GetPrivateKey()
	if err != nil {
		return nil, err
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedByte)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}
func (u *User) Check() (bool, error) {
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetHost.path("user/info"), nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", "Bearer "+u.Token)
	request.Header.Set("Content-Type", ContentTypeApplicationJSON)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return false, nil
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (u *User) Login() error {
	b, err := json.Marshal(LoginRequest{ID: u.ID})
	if err != nil {
		return err
	}
	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetHost.path("user/login"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	response := LoginResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	response.EncryptedCodeByte, err = base64.RawStdEncoding.DecodeString(response.EncryptedCode)
	if err != nil {
		return err
	}
	code, err := u.Decrypt(response.EncryptedCodeByte)
	if err != nil {
		return err
	}
	err = u.verify(string(code))
	if err != nil {
		return err
	}
	u.Store()

	return nil
}

func (u *User) verify(code string) error {
	b, err := json.Marshal(VerifyRequest{
		ID:   u.ID,
		Code: code,
	})
	if err != nil {
		return err
	}

	timeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetHost.path("user/verify"), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("somthing failed: %+v", resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	rV := VerifyResponse{}
	err = json.Unmarshal(body, &rV)
	if err != nil {
		return err
	}
	u.Token = rV.Token

	return nil
}

func (u *User) Store() {
	file, err := os.Create(u.UserFile)
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(j)
	if err != nil {
		panic(err)
	}
}
