package rabbitmq

import (
	"context"
	"fmt"
	"github.com/mohsenHa/messenger/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type rabbitmqChannel struct {
	wg                     *sync.WaitGroup
	done                   <-chan bool
	rabbitmq               *Rabbitmq
	inputChannel           chan []byte
	outputChannel          chan Message
	exchange               string
	queue                  string
	maxRetryPolicy         int
	closeSignalChannel     <-chan bool
	heartBeatSignalChannel chan<- bool
}
type rabbitmqChannelParams struct {
	rabbitmq               *Rabbitmq
	exchange               string
	queue                  string
	bufferSize             int
	maxRetryPolicy         int
	closeSignalChannel     <-chan bool
	heartBeatSignalChannel chan<- bool
}

const (
	timeForCallAgainDuration = 10
	retriesToTimeRatio       = 2
)

func newChannel(done <-chan bool, wg *sync.WaitGroup, rabbitmqChannelParams rabbitmqChannelParams) *rabbitmqChannel {
	conn := rabbitmqChannelParams.rabbitmq.connection
	WaitForConnection(rabbitmqChannelParams.rabbitmq)
	ch := openChannel(conn)
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		failOnError(err, "failed to close a channel")
	}(ch)

	err := ch.ExchangeDeclare(
		rabbitmqChannelParams.exchange, // name
		"topic",                        // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		ch := openChannel(conn)
		err = ch.ExchangeDeclarePassive(
			rabbitmqChannelParams.exchange, // name
			"topic",                        // type
			true,                           // durable
			false,                          // auto-deleted
			false,                          // internal
			false,                          // no-wait
			nil,                            // arguments
		)
		failOnError(err, "Failed to declare an exchange")
	}
	_, errQueueDeclare := ch.QueueDeclare(
		rabbitmqChannelParams.queue, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)

	if errQueueDeclare != nil {
		ch := openChannel(conn)
		_, errQueueDeclare := ch.QueueDeclarePassive(
			rabbitmqChannelParams.queue, // name
			true,                        // durable
			false,                       // delete when unused
			false,                       // exclusive
			false,                       // no-wait
			nil,                         // arguments
		)
		failOnError(errQueueDeclare, "Failed to declare a queue")
	}

	errQueueBind := ch.QueueBind(
		rabbitmqChannelParams.queue,    // queue name
		"",                             // routing key
		rabbitmqChannelParams.exchange, // exchange
		false,
		nil)
	failOnError(errQueueBind, "Failed to bind a queue")

	rc := &rabbitmqChannel{
		done:                   done,
		wg:                     wg,
		exchange:               rabbitmqChannelParams.exchange,
		queue:                  rabbitmqChannelParams.queue,
		rabbitmq:               rabbitmqChannelParams.rabbitmq,
		maxRetryPolicy:         rabbitmqChannelParams.maxRetryPolicy,
		closeSignalChannel:     rabbitmqChannelParams.closeSignalChannel,
		heartBeatSignalChannel: rabbitmqChannelParams.heartBeatSignalChannel,
		inputChannel:           make(chan []byte, rabbitmqChannelParams.bufferSize),
		outputChannel:          make(chan Message, rabbitmqChannelParams.bufferSize),
	}
	rc.start()

	return rc
}

func openChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	return ch
}

func (rc *rabbitmqChannel) GetInputChannel() chan<- []byte {
	return rc.inputChannel
}

func (rc *rabbitmqChannel) GetOutputChannel() <-chan Message {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		rc.startOutput()
	}()

	return rc.outputChannel
}

func (rc *rabbitmqChannel) start() {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		go rc.startInput()
	}()
}

func (rc *rabbitmqChannel) startOutput() {
	rc.wg.Add(1)
	WaitForConnection(rc.rabbitmq)
	go func() {
		defer rc.wg.Done()
		ch, err := rc.rabbitmq.connection.Channel()

		failOnError(err, "Failed to open a channel")
		if err != nil {
			rc.callMeNextTime(rc.startOutput)

			return
		}

		defer func(ch *amqp.Channel) {
			err = ch.Close()
			failOnError(err, "Failed to close channel")
		}(ch)

		msgs, errConsume := ch.Consume(
			rc.queue, // queue
			"",       // consumer
			false,    // auto-ack
			false,    // exclusive
			false,    // no-local
			false,    // no-wait
			nil,      // arguments
		)
		failOnError(errConsume, "failed to consume")
		if err != nil {
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
			case <-rc.closeSignalChannel:
				fmt.Println("Receive close signal")
				return
			case msg := <-msgs:
				rc.wg.Add(1)
				go func() {
					defer rc.wg.Done()

					rc.outputChannel <- Message{
						Body: msg.Body,
						Ack: func() error {
							return msg.Ack(false)
						},
					}
					rc.sendHeartBeatSignal()
				}()
			}
		}
	}()
}
func (rc *rabbitmqChannel) sendHeartBeatSignal() {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		rc.heartBeatSignalChannel <- true
	}()
}
func (rc *rabbitmqChannel) startInput() {
	rc.wg.Add(1)
	WaitForConnection(rc.rabbitmq)
	go func() {
		defer rc.wg.Done()

		ch, err := rc.rabbitmq.connection.Channel()
		failOnError(err, "Failed to open a channel")
		if err != nil {
			rc.callMeNextTime(rc.startInput)

			return
		}
		defer func(ch *amqp.Channel) {
			err = ch.Close()
			failOnError(err, "Failed to close channel")
		}(ch)

		for {
			if ch.IsClosed() {
				rc.callMeNextTime(rc.startInput)

				return
			}
			select {
			case <-rc.done:

				return
			case <-rc.closeSignalChannel:
				fmt.Println("Receive close signal")
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
		logger.Logger.Error(fmt.Sprintf("job failed after %d tries", tries))

		return
	}
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		time.Sleep(time.Second * time.Duration(tries*retriesToTimeRatio))
		errPWC := ch.PublishWithContext(context.Background(),
			rc.exchange, // exchange
			"",          // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			})
		failOnError(errPWC, "failed on ACK")
		if errPWC != nil {
			rc.publishToRabbitmq(ch, msg, tries+1)
		}
	}()
}

func (rc *rabbitmqChannel) callMeNextTime(f func()) {
	rc.wg.Add(1)
	go func() {
		time.Sleep(time.Second * timeForCallAgainDuration)
		defer rc.wg.Done()
		f()
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Logger.Error(err.Error())
		fmt.Println(err, msg)
	}
}
