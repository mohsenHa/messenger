package keygenerator

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
)

type Service struct {
	config Config
}

type Config struct {
	KeyLength uint8  `koanf:"key_length"`
	IdPrefix  string `koanf:"id_prefix"`
	IdPostfix string `koanf:"id_postfix"`
}

func New(config Config) Service {
	return Service{
		config: config,
	}
}

func (s Service) CreateCode() (string, error) {
	if s.config.KeyLength == 0 {
		return "", fmt.Errorf("key length must greater than %d", 0)
	}
	return random.String(s.config.KeyLength), nil
}

func (s Service) EncryptCode(code, publicKey string) (string, error) {
	pk, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	encryptedCode, err := encryptdecrypt.Encrypt(pk, []byte(code))
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(encryptedCode), nil
}

func (s Service) CreateUserId(publicKey string) string {
	md5Prefix := encryptdecrypt.GetMD5Hash(s.config.IdPrefix)
	md5PublicKey := encryptdecrypt.GetMD5Hash(publicKey)
	md5Postfix := encryptdecrypt.GetMD5Hash(s.config.IdPostfix)
	return encryptdecrypt.GetMD5Hash(md5Prefix + md5PublicKey + md5Postfix)
}
