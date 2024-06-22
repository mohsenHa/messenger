package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mohsenHa/messenger/repository/inmemory/inmemoryuser"

	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/delivery/httpserver"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/keygenerator"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/messagevalidator"
	"github.com/mohsenHa/messenger/validator/uservalidator"
)

func main() {
	wg := &sync.WaitGroup{}
	done := make(chan bool)

	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	_ = logger.NewLogger(cfg.Logger)

	rSvcs, rVal := setupServices(cfg, wg, done)

	server := httpserver.New(cfg, rSvcs, rVal)
	go func() {
		server.Serve()
	}()

	if cfg.Application.EnableProfiling {
		profiling(cfg, wg, done)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	close(done)

	wg.Wait()

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(cfg.Application.GracefulShutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully..")
	<-ctxWithTimeout.Done()
}

func profiling(cfg config.Config, wg *sync.WaitGroup, done <-chan bool) {
	fmt.Printf("Profiling enabled on port %d\n", cfg.Application.ProfilingPort)
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Application.ProfilingPort),
		ReadTimeout: time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-done
		timeout := 5
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
}

func setupServices(cfg config.Config, wg *sync.WaitGroup, done chan bool) (requiredServices httpserver.RequiredServices, requiredValidators httpserver.RequiredValidators) {
	userRepo := inmemoryuser.New()

	keyGen := keygenerator.New(cfg.KeyGenerator)
	authSvc := authservice.New(cfg.Auth)

	rmq := rabbitmq.New(done, wg, cfg.Rabbitmq)

	requiredServices.MessageService = messageservice.New(rmq, userRepo)
	requiredServices.UserService = userservice.New(userRepo, authSvc, keyGen)
	requiredServices.AuthService = authSvc

	requiredValidators.UserValidator = uservalidator.New(userRepo, keyGen)
	requiredValidators.MessageValidator = messagevalidator.New(userRepo, keyGen)

	return requiredServices, requiredValidators
}
