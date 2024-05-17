package config

import (
	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/repository/mysql"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/keygenerator"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application  Application         `koanf:"application"`
	HTTPServer   HTTPServer          `koanf:"http_server"`
	Mysql        mysql.Config        `koanf:"mysql"`
	Rabbitmq     rabbitmq.Config     `koanf:"rabbitmq"`
	KeyGenerator keygenerator.Config `koanf:"key_generator"`
	Auth         authservice.Config  `koanf:"auth"`
}
