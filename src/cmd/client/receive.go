package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mohsenHa/messenger/param/messageparam"
)

func Receive(wg *sync.WaitGroup, done <-chan bool, user User) {
	url := targetHost.ws("message/receive?token=" + user.Token)
	fmt.Println(url)
	c, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	bufferSize := 100
	messageChannel := make(chan messageparam.SendMessage, bufferSize)
	closedChannel := make(chan bool)
	go startListen(c, messageChannel, closedChannel)
	wg.Add(1)
	go func() {
		defer func(c *websocket.Conn) {
			wg.Done()
			err = c.Close()
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
				switch msg.Type {
				case messageparam.SendMessageDeliverType:
					deliverReceived(msg)
				case messageparam.SendMessageMessageType:
					messageReceived(msg, user)
				default:
					fmt.Printf("Message received with invalid type %v", msg)
				}
			}
		}
	}()
}

func deliverReceived(msg messageparam.SendMessage) {
	fmt.Printf("Message %s delivered to: %v\t%v\n", msg.ID, msg.From.ID,
		msg.SendTime.Format("2006-01-02 15:04:05"))
}

func messageReceived(msg messageparam.SendMessage, user User) {
	decodeString, err := base64.RawStdEncoding.DecodeString(msg.Body)
	if err != nil {
		fmt.Println("error on decode message", err)

		return
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(nil, user.PrivateKeyRSA, decodeString)
	if err != nil {
		fmt.Println("error on decrypt message", err)

		return
	}
	fmt.Printf("Message From: %v\t%v\n", msg.From.ID, msg.SendTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Message: %s\n", decryptedBytes)
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
