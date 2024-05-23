package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mohsenHa/messenger/param/messageparam"
	"log"
	"sync"
	"time"
)

func Receive(wg *sync.WaitGroup, done <-chan bool, user User) {
	url := targetHost.ws("message/receive?token=" + user.Token)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}

	messageChannel := make(chan messageparam.SendMessage, 100)
	closedChannel := make(chan bool)
	go startListen(c, messageChannel, closedChannel)
	wg.Add(1)
	go func() {
		defer func(c *websocket.Conn) {
			wg.Done()
			err := c.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(c)

		fmt.Println("Start listening")
		for {
			select {
			case <-closedChannel:
				fmt.Println("End listening")
				time.Sleep(time.Second)
				go Receive(wg, done, user)
				return
			case <-done:
				fmt.Println("End listening")
				return
			case msg := <-messageChannel:
				decodeString, err := base64.RawStdEncoding.DecodeString(msg.Body)
				if err != nil {
					fmt.Println("error on decode message", err)
					continue
				}
				decryptedBytes, err := rsa.DecryptPKCS1v15(nil, user.PrivateKeyRSA, decodeString)
				if err != nil {
					fmt.Println("error on decrypt message", err)
					continue
				}
				fmt.Printf("Message From: %v\n", msg.From.Id)
				fmt.Printf("Message: %s\n", decryptedBytes)
			}
		}
	}()
}

func startListen(c *websocket.Conn, messageChannel chan<- messageparam.SendMessage, closedChannel chan bool) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			close(closedChannel)
			return
		}
		m := messageparam.SendMessage{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			log.Println(err, message)
			continue
		}
		messageChannel <- m
	}

}
