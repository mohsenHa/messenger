package main

import (
	"bufio"
	"fmt"
	"github.com/mohsenHa/messenger/pkg/encryptdecrypt"
	"os"
)

func Messenger(user User) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("For send message please enter target id")
	for scanner.Scan() {

		to := scanner.Text()

		if to == "exit" {
			return
		}
		if to == "id" {
			fmt.Printf("Your ID is %s \n", user.ID)

			continue
		}
		publicKey, err := GetPublicKey(GetPublicKeyRequest{
			ID:    to,
			Token: user.Token,
		})
		if err != nil {
			fmt.Println(err)

			continue
		}
		fmt.Printf("You send message to %s for end chat enter exit \n", to)
		startSendMessage(publicKey.PublicKey, to, user.Token)
		fmt.Println("For send message please enter target id")
	}
}

func startSendMessage(publicKey, to, token string) {
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
			ToID:    to,
			Token:   token,
		})
		if err != nil {
			fmt.Println(err)

			continue
		}
		fmt.Printf("Message sent with id %s\n", send.ID)
	}
}
