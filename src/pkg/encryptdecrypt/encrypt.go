package encryptdecrypt

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func Encrypt(pk string, message string) (string, error) {
	const op = "encryptdecrypt.Encrypt"
	publicKeyPEM, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return "", richerror.New(op).WithErr(fmt.Errorf("invalid public key")).
			WithMessage(errmsg.ErrorMsgInvalidInput).WithKind(richerror.KindInvalid)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		return "", richerror.New(op).WithErr(fmt.Errorf("invalid public key")).
			WithMessage(errmsg.ErrorMsgInvalidInput).WithKind(richerror.KindInvalid)

	}
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return "", richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	encryptedCode, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(message))

	return base64.RawStdEncoding.EncodeToString(encryptedCode), nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
