package main

import (
	"context"
	"fmt"
	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/delivery/httpserver"
	"github.com/mohsenHa/messenger/repository/migrator/mysqlmigrator"
	"github.com/mohsenHa/messenger/repository/mysql"
	"github.com/mohsenHa/messenger/repository/mysql/mysqluser"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/keygenerator"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/messagevalidator"
	"github.com/mohsenHa/messenger/validator/uservalidator"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	done := make(chan bool)

	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	mgr := mysqlmigrator.New(cfg.Mysql)
	mgr.Up()

	rSvcs, rVal := setupServices(cfg, &wg, done)

	server := httpserver.New(cfg, rSvcs, rVal)
	go func() {
		server.Serve()
	}()

	if cfg.Application.EnableProfiling {
		wg.Add(1)
		go func() {
			defer wg.Done()
			http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.ProfilingPort), nil)
		}()
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

func setupServices(cfg config.Config, wg *sync.WaitGroup, done chan bool) (requiredServices httpserver.RequiredServices, requiredValidators httpserver.RequiredValidators) {
	mysqlRepo := mysql.New(cfg.Mysql)
	userRepo := mysqluser.New(mysqlRepo)
	keyGen := keygenerator.New(cfg.KeyGenerator)
	authSvc := authservice.New(cfg.Auth)

	rmq := rabbitmq.New(done, wg, cfg.Rabbitmq)

	requiredServices.MessageService = messageservice.New(rmq, userRepo)
	requiredServices.UserService = userservice.New(userRepo, authSvc, keyGen)
	requiredServices.AuthService = authSvc

	requiredValidators.UserValidator = uservalidator.New(userRepo, keyGen)
	requiredValidators.MessageValidator = messagevalidator.New(userRepo, keyGen)
	return
}
