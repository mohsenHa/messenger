package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
)

type hostType string

func (ht hostType) path(api string) string {
	return "http://" + string(ht) + "/" + api
}
func (ht hostType) ws(api string) string {
	return "ws://" + string(ht) + "/" + api
}

const targetHost hostType = "messenger.local"

func main() {
	wg := &sync.WaitGroup{}
	done := make(chan bool)
	user := loadUser()

	fmt.Printf("Your Id is %s\n", user.Id)

	Receive(wg, done, user)
	Messenger(user)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Graceful shutdown")
	close(done)
	fmt.Println("Wait for done")
	wg.Wait()

}

func loadUser() User {
	args := os.Args[1:]

	userFile := "user.json"
	if len(args) > 0 {
		userFile = args[0] + ".json"
	}

	user := User{}
	file, err := os.Open(userFile)
	if err != nil {
		return createUser(userFile)
	}

	j, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(j, &user)
	if err != nil {
		panic(err)
	}
	keyPEMbyte, err := base64.StdEncoding.DecodeString(user.PrivateKey)
	if err != nil {
		panic(err)
	}
	keyPEM, _ := pem.Decode(keyPEMbyte)
	key, err := x509.ParsePKCS1PrivateKey(keyPEM.Bytes)
	if err != nil {
		panic(err)
	}
	user.PrivateKeyRSA = key

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
	return user
}

func createUser(userFile string) User {
	user, err := NewUser()
	if err != nil {
		panic(err)
	}
	file, err := os.Create(userFile)
	if err != nil {
		panic(err)
	}
	j, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(j)
	if err != nil {
		panic(err)
	}
	return user
}
