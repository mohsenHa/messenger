package main

import (
	"github.com/mohsenHa/messenger/config"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
)

func main() {
	cfg := config.Load("config.yml")

	logger.NewLogger(cfg.Logger)
	logger.L().Debug(loggerentity.CategoryRabbitMQ, loggerentity.SubCategoryRabbitMQConnection,
		"rabbitmq connection error", map[loggerentity.ExtraKey]interface{}{
			loggerentity.ExtraKeyErrorMessage: "test",
		})
	logger.NewLog("rabbitmq connection error").
		WithCategory(loggerentity.CategoryRabbitMQ).
		WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
		With(loggerentity.ExtraKeyErrorMessage, "test").
		Debug()

}
