package main

import (
	"context"
	"fmt"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/delivery/httpserver"
	"github.com/mohsenHa/messenger/repository/migrator/mysqlmigrator"
	"github.com/mohsenHa/messenger/repository/mysql"
	"github.com/mohsenHa/messenger/repository/mysql/mysqluser"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/uservalidator"
	"os"
	"os/signal"
)

const (
	JwtSignKey = ""
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	mgr := mysqlmigrator.New(cfg.Mysql)
	mgr.Down()
	mgr.Up()

	rSvcs, rVal := setupServices(cfg)

	server := httpserver.New(cfg, rSvcs, rVal)
	go func() {
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully..")
	<-ctxWithTimeout.Done()
}

func setupServices(cfg config.Config) (requiredServices httpserver.RequiredServices, requiredValidators httpserver.RequiredValidators) {
	mysqlRepo := mysql.New(cfg.Mysql)

	userRepo := mysqluser.New(mysqlRepo)

	requiredServices.MessageService = messageservice.New(cfg.Rabbitmq)
	requiredServices.UserService = userservice.New(userRepo, cfg.UserService)

	requiredValidators.UserValidator = uservalidator.New(userRepo)

	return
}
