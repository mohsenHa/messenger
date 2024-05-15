package encryptdecrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"github.com/mohsenHa/messenger/pkg/richerror"
)

func Encrypt(pk string, message []byte) ([]byte, error) {
	const op = "encryptdecrypt.Encrypt"
	publicKeyPEM := []byte(base64.StdEncoding.EncodeToString([]byte(pk)))
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), message)
}
