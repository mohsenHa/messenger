package rabbitmq

import (
	"context"
	"fmt"
	"github.com/mohsenHa/messenger/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type OutputChannel struct {
	closeChannel         <-chan bool
	internalCloseChannel chan bool
	messageChannel       chan Message
	isClosed             bool
}
type rabbitmqChannel struct {
	wg                     *sync.WaitGroup
	done                   <-chan bool
	rabbitmq               *Rabbitmq
	inputChannel           chan []byte
	outputChannels         []*OutputChannel
	outputChannelMutex     *sync.RWMutex
	exchange               string
	queue                  string
	maxRetryPolicy         int
	closeSignalChannel     <-chan bool
	heartBeatSignalChannel chan<- bool
	openOutputConsumer     *sync.Once
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
		true,                        // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)

	if errQueueDeclare != nil {
		ch := openChannel(conn)
		_, errQueueDeclare := ch.QueueDeclarePassive(
			rabbitmqChannelParams.queue, // name
			true,                        // durable
			true,                        // delete when unused
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
		outputChannels:         make([]*OutputChannel, rabbitmqChannelParams.bufferSize),
		outputChannelMutex:     &sync.RWMutex{},
		openOutputConsumer:     &sync.Once{},
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

func (rc *rabbitmqChannel) GetOutputChannel(closeChanelSignal <-chan bool) <-chan Message {
	rc.openOutputConsumer.Do(func() {
		rc.wg.Add(1)
		go func() {
			defer rc.wg.Done()
			rc.startOutput()
		}()
	})

	return rc.NewOutputChannel(closeChanelSignal)
}

func (rc *rabbitmqChannel) NewOutputChannel(closeChanelSignal <-chan bool) <-chan Message {
	channel := make(chan Message)
	outputChanel := &OutputChannel{
		closeChannel:         closeChanelSignal,
		internalCloseChannel: make(chan bool),
		messageChannel:       channel,
		isClosed:             false,
	}
	rc.outputChannelMutex.Lock()
	rc.outputChannels = append(rc.outputChannels, outputChanel)
	rc.outputChannelMutex.Unlock()
	rc.WaitForCloseOutputChannel(outputChanel)

	return channel
}

func (rc *rabbitmqChannel) WaitForCloseOutputChannel(outputChanel *OutputChannel) {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		for {
			select {
			case <-rc.done:
				return
			case <-outputChanel.closeChannel:
				outputChanel.isClosed = true
				return
			case <-outputChanel.internalCloseChannel:
				outputChanel.isClosed = true
				return
			}
		}
	}()
}

func (rc *rabbitmqChannel) GetHeartbeatChannel() chan<- bool {
	return rc.heartBeatSignalChannel
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

		for index, c := range rc.outputChannels {
			if c.isClosed {
				rc.deleteOutputChannel(index)
				continue
			}
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
				fmt.Println("Timeout on all ACKs")
				return
			default:
				time.Sleep(time.Second)
			}
			if ackCount == 0 {
				fmt.Println("All ack is done")
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
	rc.outputChannels = append(rc.outputChannels[:index], rc.outputChannels[index+1:]...)
	rc.outputChannelMutex.Unlock()
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
