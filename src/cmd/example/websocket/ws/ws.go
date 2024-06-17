package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func producer(remoteAddr string, channel chan string) {
	sleep := 5
	for {
		channel <- remoteAddr
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func main() {
	err := http.ListenAndServe(":8040", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
			panic(err)
		}

		defer func(conn net.Conn) {
			err = conn.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(conn)

		done := make(chan bool)
		go readMessage(conn, done)

		channel := make(chan string)
		go producer(r.RemoteAddr, channel)
		go writeMessage(conn, channel)

		<-done
	}))
	if err != nil {
		panic(err)
	}
}

func readMessage(conn io.ReadWriter, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {

			log.Print(err)
			done <- true

			return
		}

		fmt.Println(string(msg))

		fmt.Println("opCode", opCode)
	}
}

func writeMessage(conn io.Writer, channel <-chan string) {
	for data := range channel {
		err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(data))
		if err != nil {
			panic(err)
		}
	}
}
