package config

import (
	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/repository/mysql"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/keygenerator"
)

type Application struct {
	GracefulShutdownTimeout int  `koanf:"graceful_shutdown_timeout"`
	EnableProfiling         bool `koanf:"enable_profiling"`
	ProfilingPort           int  `koanf:"profiling_port"`
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
