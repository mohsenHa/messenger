package encryptdecrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func Encrypt(pk, message string) (string, error) {
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
	if err != nil {
		return "", richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return base64.RawStdEncoding.EncodeToString(encryptedCode), nil
}

func GetMD5Hash(text string) string {
	return GetHash(text)
}

func GetHash(text string) string {
	sh := sha256.New()
	sh.Write([]byte(text))
	bs := sh.Sum(nil)

	return hex.EncodeToString(bs)
}
