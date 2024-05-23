package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type User struct {
	Id            string          `json:"id"`
	Token         string          `json:"token"`
	PublicKey     string          `json:"public_key"`
	PrivateKey    string          `json:"private_key"`
	PrivateKeyRSA *rsa.PrivateKey `json:"-"`
}

func (u User) GetPrivateKey() (*rsa.PrivateKey, error) {
	keyPEM, err := base64.StdEncoding.DecodeString(u.PrivateKey)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyPEM)

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privatekey, nil
}

func NewUser() (User, error) {
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
	decryptedBytes, err := rsa.DecryptPKCS1v15(nil, privateKey, rR.EncryptedCodeByte)
	if err != nil {
		return User{}, err
	}

	vR, err := Verify(VerifyRequest{
		Id:   rR.Id,
		Code: string(decryptedBytes),
	})
	if err != nil {
		return User{}, err
	}

	return User{
		Id:            vR.Id,
		Token:         vR.Token,
		PublicKey:     publicKeyBase64,
		PrivateKey:    privateKeyBase64,
		PrivateKeyRSA: privateKey,
	}, nil

}
