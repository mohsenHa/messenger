package keygenerator

import (
	"fmt"

	"github.com/labstack/gommon/random"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
)

type Service struct {
	config Config
}

type Config struct {
	KeyLength uint8  `koanf:"key_length"`
	IDPrefix  string `koanf:"id_prefix"`
	IDPostfix string `koanf:"id_postfix"`
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
	return encryptdecrypt.Encrypt(publicKey, code)
}

func (s Service) CreateUserID(publicKey string) string {
	md5Prefix := encryptdecrypt.GetMD5Hash(s.config.IDPrefix)
	md5PublicKey := encryptdecrypt.GetMD5Hash(publicKey)
	md5Postfix := encryptdecrypt.GetMD5Hash(s.config.IDPostfix)

	return encryptdecrypt.GetMD5Hash(md5Prefix + md5PublicKey + md5Postfix)
}
