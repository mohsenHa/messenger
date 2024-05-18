package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/delivery/httpserver/messagehandler"
	"github.com/mohsenHa/messenger/delivery/httpserver/userhandler"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/messagevalidator"
	"github.com/mohsenHa/messenger/validator/uservalidator"
	"go.uber.org/zap"
)

type Server struct {
	config         config.Config
	Router         *echo.Echo
	messageHandler messagehandler.Handler
	userHandler    userhandler.Handler
}

type RequiredServices struct {
	UserService    userservice.Service
	MessageService messageservice.Service
	AuthService    authservice.Service
}

type RequiredValidators struct {
	UserValidator    uservalidator.Validator
	MessageValidator messagevalidator.Validator
}

func New(config config.Config, services RequiredServices, validators RequiredValidators) Server {
	return Server{
		Router:         echo.New(),
		config:         config,
		messageHandler: messagehandler.New(services.MessageService, services.AuthService, validators.MessageValidator),
		userHandler:    userhandler.New(services.UserService, validators.UserValidator),
	}
}

func (s Server) Serve() {
	// Middleware

	s.Router.Use(middleware.RequestID())

	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			logger.Logger.Named("http-server").Info("request",
				zap.String("request_id", v.RequestID),
				zap.String("host", v.Host),
				zap.String("content-length", v.ContentLength),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.String("error", errMsg),
				zap.String("remote_ip", v.RemoteIP),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	//s.Router.Use(middleware.Logger())

	s.Router.Use(middleware.Recover())

	// Routes
	s.Router.GET("/health-check", s.healthCheck)

	s.messageHandler.SetRoutes(s.Router.Group("message"))

	s.userHandler.SetRoutes(s.Router.Group("user"))

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
