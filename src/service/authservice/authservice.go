package authservice

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenHa/messenger/entity"
	"os"
	"strings"
	"time"
)

type Config struct {
	PublicKeyPath         string        `koanf:"public_key_path"`
	PrivateKeyPath        string        `koanf:"private_key_path"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject         string        `koanf:"access_subject"`
	RefreshSubject        string        `koanf:"refresh_subject"`
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
	return s.createToken(user.Id, s.config.AccessExpirationTime)
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		keyData, _ := os.ReadFile(s.config.PublicKeyPath)
		key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) createToken(id string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		Id: id,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	keyDataPrivate, _ := os.ReadFile(s.config.PrivateKeyPath)
	keyprivate, _ := jwt.ParseRSAPrivateKeyFromPEM(keyDataPrivate)
	tokenString, err := accessToken.SignedString(keyprivate)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
