package rabbitmq

import (
	"fmt"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func (rc *rabbitmqChannel) GetOutputChannel(outputChannelCloseSignal <-chan bool) <-chan Message {
	rc.openOutputConsumer.Do(func() {
		rc.wg.Add(1)
		go func() {
			defer rc.wg.Done()
			rc.startOutput()
		}()
	})

	return rc.NewOutputChannel(outputChannelCloseSignal)
}

func (rc *rabbitmqChannel) NewOutputChannel(outputChannelCloseSignal <-chan bool) <-chan Message {
	channel := make(chan Message)
	outputChanel := &OutputChannel{
		outputChannelCloseSignal: outputChannelCloseSignal,
		internalCloseChannel:     make(chan bool),
		messageChannel:           channel,
		isClosed:                 false,
	}

	index := rc.addOutputChannel(outputChanel)
	rc.WaitForCloseOutputChannel(outputChanel, index)

	return channel
}

func (rc *rabbitmqChannel) addOutputChannel(outputChanel *OutputChannel) int {
	rc.outputChannelMutex.Lock()
	index := len(rc.outputChannels)
	rc.outputChannels = append(rc.outputChannels, outputChanel)
	rc.outputChannelMutex.Unlock()
	return index
}

func (rc *rabbitmqChannel) WaitForCloseOutputChannel(outputChanel *OutputChannel, index int) {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		channelDeleted := "Output channel is deleted"
		select {
		case <-rc.done:
			return
		case <-outputChanel.outputChannelCloseSignal:
			logger.L().Debugf(channelDeleted)
			rc.deleteOutputChannel(index)
			return
		case <-outputChanel.internalCloseChannel:
			logger.L().Debugf(channelDeleted)
			rc.deleteOutputChannel(index)
			return
		case <-rc.rabbitMQCloseSignal:
			logger.L().Debugf(channelDeleted)
			rc.deleteOutputChannel(index)
			return
		}
	}()
}

func (rc *rabbitmqChannel) GetHeartbeatChannel() chan<- bool {
	return rc.heartBeatSignalChannel
}

func (rc *rabbitmqChannel) startOutput() {
	WaitForConnection(rc.rabbitmq)
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		ch, err := rc.rabbitmq.connection.Channel()

		if err != nil {
			logger.NewLog("Failed to open a channel").
				WithCategory(loggerentity.CategoryRabbitMQ).
				WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
			rc.callMeNextTime(rc.startOutput)

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

		msgs, err := ch.Consume(
			rc.queue, // queue
			"",       // consumer
			false,    // auto-ack
			false,    // exclusive
			false,    // no-local
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			logger.NewLog("Failed to start consume").
				WithCategory(loggerentity.CategoryRabbitMQ).
				WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()

			rc.callMeNextTime(rc.startOutput)

			return
		}

		for {
			if ch.IsClosed() {
				rc.callMeNextTime(rc.startOutput)

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
				for _, c := range rc.outputChannels {
					close(c.internalCloseChannel)
				}
				return
			case msg := <-msgs:
				rc.receivedMessage(msg)
			}
		}
	}()
}

func (rc *rabbitmqChannel) receivedMessage(msg amqp.Delivery) {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()

		ackChan := make(chan bool)
		ackCount := 0
		message := Message{
			Body: msg.Body,
			Ack: func() error {
				ackChan <- true
				return nil
			},
		}

		for _, c := range rc.outputChannels {
			ackCount++
			c.messageChannel <- message
		}
		rc.processSubAck(ackChan, ackCount, msg)
		rc.sendHeartBeatSignal()
	}()

}

func (rc *rabbitmqChannel) processSubAck(ackChan <-chan bool, ackCount int, msg amqp.Delivery) {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		askTimeout := time.After(time.Second * 5)
		for {
			select {
			case <-ackChan:
				ackCount--
			case <-askTimeout:
				logger.NewLog("Timeout on all ACKs").
					WithCategory(loggerentity.CategoryRabbitMQ).
					WithSubCategory(loggerentity.SubCategoryRabbitMQChannel).
					Debug()
				return
			default:
				time.Sleep(time.Second)
			}
			if ackCount == 0 {
				logger.NewLog("All Ack is done").
					WithCategory(loggerentity.CategoryRabbitMQ).
					WithSubCategory(loggerentity.SubCategoryRabbitMQAck).
					Debug()
				err := msg.Ack(false)
				if err != nil {
					fmt.Println("error in ack", err)
				}
				return
			}
		}
	}()
}

func (rc *rabbitmqChannel) deleteOutputChannel(index int) {
	rc.outputChannelMutex.Lock()
	close(rc.outputChannels[index].messageChannel)
	rc.outputChannels = append(rc.outputChannels[:index], rc.outputChannels[index+1:]...)
	rc.outputChannelMutex.Unlock()
	if len(rc.outputChannels) == 0 {
		close(rc.noOutputConsumer)
	}
}

func (rc *rabbitmqChannel) sendHeartBeatSignal() {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		rc.heartBeatSignalChannel <- true
	}()
}
