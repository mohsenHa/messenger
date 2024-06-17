package authservice

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenHa/messenger/entity"
	"github.com/mohsenHa/messenger/pkg/errmsg"
	"os"
	"strings"
	"time"
)

type Config struct {
	PublicKeyPath   string `koanf:"public_key_path"`
	PrivateKeyPath  string `koanf:"private_key_path"`
	AccessTTLSecond int    `koanf:"access_ttl_second"`
	AccessSubject   string `koanf:"access_subject"`
	RefreshSubject  string `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, time.Second*time.Duration(s.config.AccessTTLSecond))
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(_ *jwt.Token) (interface{}, error) {
		return s.GetPublicKey()
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf(errmsg.ErrorMsgUnauthorized)
}

func (s Service) GetPublicKey() (interface{}, error) {
	keyData, err := os.ReadFile(s.config.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s Service) GetPrivateKey() (interface{}, error) {
	keyDataPrivate, err := os.ReadFile(s.config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyDataPrivate)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s Service) GetSigningMethod() *jwt.SigningMethodRSA {
	return jwt.SigningMethodRS512
}

func (s Service) createToken(id string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		ID: id,
	}
	accessToken := jwt.NewWithClaims(s.GetSigningMethod(), claims)
	keyprivate, err := s.GetPrivateKey()
	if err != nil {
		return "", err
	}
	tokenString, err := accessToken.SignedString(keyprivate)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
