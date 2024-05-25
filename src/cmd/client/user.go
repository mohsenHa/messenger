package main

import (
	"bytes"
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
)

type LoginRequest struct {
	Id string `json:"id"`
}
type LoginResponse struct {
	EncryptedCode     string `json:"encrypted_code"`
	EncryptedCodeByte []byte
}

type VerifyRequest struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}
type VerifyResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type User struct {
	Id            string          `json:"id"`
	Token         string          `json:"token"`
	PublicKey     string          `json:"public_key"`
	PrivateKey    string          `json:"private_key"`
	PrivateKeyRSA *rsa.PrivateKey `json:"-"`
	UserFile      string          `json:"-"`
}

const ContentTypeApplicationJson = "application/json"

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
	bytes, err := x509.MarshalPKIXPublicKey(pub.(*rsa.PublicKey))
	if err != nil {
		return User{}, err
	}
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
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
		Id:            rR.Id,
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
	request, err := http.NewRequest(http.MethodGet, targetHost.path("user/info"), nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", "Bearer "+u.Token)
	request.Header.Set("Content-Type", ContentTypeApplicationJson)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return false, nil
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func (u *User) Login() error {
	b, err := json.Marshal(LoginRequest{Id: u.Id})
	if err != nil {
		return err
	}

	resp, err := http.Post(targetHost.path("user/login"), ContentTypeApplicationJson, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

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
		Id:   u.Id,
		Code: code,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(targetHost.path("user/verify"), ContentTypeApplicationJson, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("somthing failed: %+v", resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
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
