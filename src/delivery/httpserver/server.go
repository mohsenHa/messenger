package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/delivery/httpserver/messagehandler"
	"github.com/mohsenHa/messenger/delivery/httpserver/userhandler"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"github.com/mohsenHa/messenger/service/authservice"
	"github.com/mohsenHa/messenger/service/messageservice"
	"github.com/mohsenHa/messenger/service/userservice"
	"github.com/mohsenHa/messenger/validator/messagevalidator"
	"github.com/mohsenHa/messenger/validator/uservalidator"
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
		userHandler:    userhandler.New(services.UserService, validators.UserValidator, services.AuthService),
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

			logger.NewLog("http-server").
				WithCategory(loggerentity.CategoryRequestResponse).
				WithSubCategory(loggerentity.SubCategoryInternalRequest).
				With(loggerentity.ExtraKeyRequestID, v.RequestID).
				With(loggerentity.ExtraKeyHost, v.Host).
				With(loggerentity.ExtraKeyContentLength, v.ContentLength).
				With(loggerentity.ExtraKeyProtocol, v.Protocol).
				With(loggerentity.ExtraKeyMethod, v.Method).
				With(loggerentity.ExtraKeyLatency, v.Latency).
				With(loggerentity.ExtraKeyErrorMessage, errMsg).
				With(loggerentity.ExtraKeyRemoteIP, v.RemoteIP).
				With(loggerentity.ExtraKeyResponseSize, v.ResponseSize).
				With(loggerentity.ExtraKeyURI, v.URI).
				With(loggerentity.ExtraKeyURIPath, v.URIPath).
				With(loggerentity.ExtraKeyStatusCode, v.Status).
				Info()

			return nil
		},
	}))

	//s.Router.Use(middleware.Logger())
	//
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
