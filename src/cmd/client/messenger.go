package main

import (
	"bufio"
	"fmt"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
	"os"
)

func Messenger(user User) {
	for {
		fmt.Println("For send message please enter target id")
		to := ""
		_, err := fmt.Scanln(&to)
		if err != nil {
			fmt.Println(err)
			continue
		}
		publicKey, err := GetPublicKey(GetPublicKeyRequest{
			Id:    to,
			Token: user.Token,
		})
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("You send message to %s for end chat enter exit \n", to)
		startSendMessage(publicKey.PublicKey, to, user.Token)

	}

}

func startSendMessage(publicKey string, to string, token string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "exit" {
			return
		}
		encryptedMessage, err := encryptdecrypt.Encrypt(publicKey, message)
		if err != nil {
			fmt.Println(err)
			continue
		}
		send, err := Send(SendRequest{
			Message: encryptedMessage,
			ToId:    to,
			Token:   token,
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Message sent with id %s\n", send.Id)
	}
}
