package rabbitmq

import (
	"context"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func (rc *rabbitmqChannel) GetInputChannel() chan<- []byte {
	return rc.inputChannel
}

func (rc *rabbitmqChannel) startInput() {
	rc.wg.Add(1)
	WaitForConnection(rc.rabbitmq)
	go func() {
		defer rc.wg.Done()

		ch, err := rc.rabbitmq.connection.Channel()
		if err != nil {
			rc.callMeNextTime(rc.startInput)
			logger.NewLog("Failed to open a channel").
				WithCategory(loggerentity.CategoryRabbitMQ).
				WithSubCategory(loggerentity.SubCategoryRabbitMQChannel).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
			return
		}
		defer func(ch *amqp.Channel) {
			err = ch.Close()
			if err != nil {
				logger.NewLog("Failed to close channel").
					WithCategory(loggerentity.CategoryRabbitMQ).
					WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
					With(loggerentity.ExtraKeyErrorMessage, err.Error()).
					Error()
			}
		}(ch)

		for {
			if ch.IsClosed() {
				rc.callMeNextTime(rc.startInput)

				return
			}
			select {
			case <-rc.done:

				return
			case <-rc.rabbitMQCloseSignal:
				logger.NewLog("Receive close signal").
					WithCategory(loggerentity.CategoryRabbitMQ).
					WithSubCategory(loggerentity.SubCategoryRabbitMQChannel).
					Debug()
				return

			case msg := <-rc.inputChannel:
				rc.publishToRabbitmq(ch, msg, 0)
				rc.sendHeartBeatSignal()
			}
		}
	}()
}

func (rc *rabbitmqChannel) publishToRabbitmq(ch *amqp.Channel, msg []byte, tries int) {
	if tries > rc.maxRetryPolicy {
		logger.L().Errorf("job failed after %d tries", tries)

		return
	}
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		time.Sleep(time.Second * time.Duration(tries*retriesToTimeRatio))
		err := ch.PublishWithContext(context.Background(),
			rc.exchange, // exchange
			"",          // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			})
		if err != nil {
			logger.NewLog("failed on ACK").
				WithCategory(loggerentity.CategoryRabbitMQ).
				WithSubCategory(loggerentity.SubCategoryRabbitMQAck).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
			rc.publishToRabbitmq(ch, msg, tries+1)
		}
	}()
}
