package main

import (
	"fmt"
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

var targetHost hostType = "127.0.0.1:8080"

func main() {
	wg := &sync.WaitGroup{}
	done := make(chan bool)
	args := os.Args[1:]
	userFile := "user.json"
	if len(args) > 0 {
		userFile = args[0] + ".json"
	}
	if len(args) > 1 {
		targetHost = hostType(args[1])
	}

	user, err := NewUserFromFile(userFile)
	if err != nil {
		panic(err)
	}

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
