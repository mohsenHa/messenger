package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenHa/messenger/service/authservice"
	"os"
	"time"
)

func main() {
	ttl := time.Duration(10)
	claims := authservice.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * ttl)),
		},
		ID: "1234",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	keyDataPrivate, err := os.ReadFile("./key/key")
	if err != nil {
		panic(err)
	}
	keyprivate, err := jwt.ParseRSAPrivateKeyFromPEM(keyDataPrivate)
	if err != nil {
		panic(err)
	}
	tokenString, err := accessToken.SignedString(keyprivate)
	if err != nil {
		panic(err)
	}
	fmt.Println(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &authservice.Claims{}, func(_ *jwt.Token) (interface{}, error) {
		keyData, err := os.ReadFile("./key/key.pub")
		if err != nil {
			panic(err)
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
		if err != nil {
			panic(err)
		}

		return key, nil
	})
	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(*authservice.Claims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		panic("Token is invalid")
	}

}
