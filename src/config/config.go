package config

import (
	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/repository/mysql"
	"github.com/mohsenHa/messenger/service/userservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application Application        `koanf:"application"`
	HTTPServer  HTTPServer         `koanf:"http_server"`
	Mysql       mysql.Config       `koanf:"mysql"`
	Rabbitmq    rabbitmq.Config    `koanf:"rabbitmq"`
	UserService userservice.Config `koanf:"user_service"`
}
