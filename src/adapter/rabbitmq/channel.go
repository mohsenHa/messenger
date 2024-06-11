package rabbitmq

import (
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type OutputChannel struct {
	outputChannelCloseSignal <-chan bool
	internalCloseChannel     chan bool
	messageChannel           chan Message
	isClosed                 bool
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
	rabbitMQCloseSignal    <-chan bool
	noOutputConsumer       chan<- bool
	heartBeatSignalChannel chan<- bool
	openOutputConsumer     *sync.Once
}
type rabbitmqChannelParams struct {
	rabbitmq               *Rabbitmq
	exchange               string
	queue                  string
	bufferSize             int
	maxRetryPolicy         int
	rabbitMQCloseSignal    <-chan bool
	noOutputConsumer       chan<- bool
	heartBeatSignalChannel chan<- bool
}

const (
	loggerGroupName          = "adapter.rabbitmq"
	timeForCallAgainDuration = 10
	retriesToTimeRatio       = 2
)

func newChannel(done <-chan bool, wg *sync.WaitGroup, rabbitmqChannelParams rabbitmqChannelParams) (*rabbitmqChannel, error) {
	conn := rabbitmqChannelParams.rabbitmq.connection
	WaitForConnection(rabbitmqChannelParams.rabbitmq)
	ch, err := openChannel(conn)
	if err != nil {
		return nil, err
	}

	defer func(ch *amqp.Channel) {
		err = ch.Close()
		if err != nil {
			logger.NewLog("rabbitmq close connection error").
				WithCategory(loggerentity.CategoryRabbitMQ).
				WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()
		}
	}(ch)

	err = ch.ExchangeDeclare(
		rabbitmqChannelParams.exchange, // name
		"topic",                        // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		ch, err = openChannel(conn)
		if err != nil {
			return nil, err
		}
		err = ch.ExchangeDeclarePassive(
			rabbitmqChannelParams.exchange, // name
			"topic",                        // type
			true,                           // durable
			false,                          // auto-deleted
			false,                          // internal
			false,                          // no-wait
			nil,                            // arguments
		)
		if err != nil {
			return nil, err
		}
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
		ch, err = openChannel(conn)
		if err != nil {
			return nil, err
		}
		_, errQueueDeclare = ch.QueueDeclarePassive(
			rabbitmqChannelParams.queue, // name
			true,                        // durable
			false,                       // delete when unused
			false,                       // exclusive
			false,                       // no-wait
			nil,                         // arguments
		)
		if err != nil {
			return nil, err
		}
	}

	err = ch.QueueBind(
		rabbitmqChannelParams.queue,    // queue name
		"",                             // routing key
		rabbitmqChannelParams.exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}

	rc := &rabbitmqChannel{
		done:                   done,
		wg:                     wg,
		exchange:               rabbitmqChannelParams.exchange,
		queue:                  rabbitmqChannelParams.queue,
		rabbitmq:               rabbitmqChannelParams.rabbitmq,
		maxRetryPolicy:         rabbitmqChannelParams.maxRetryPolicy,
		rabbitMQCloseSignal:    rabbitmqChannelParams.rabbitMQCloseSignal,
		noOutputConsumer:       rabbitmqChannelParams.noOutputConsumer,
		heartBeatSignalChannel: rabbitmqChannelParams.heartBeatSignalChannel,
		inputChannel:           make(chan []byte, rabbitmqChannelParams.bufferSize),
		outputChannels:         make([]*OutputChannel, rabbitmqChannelParams.bufferSize),
		outputChannelMutex:     &sync.RWMutex{},
		openOutputConsumer:     &sync.Once{},
	}
	rc.start()

	return rc, nil
}

func (rc *rabbitmqChannel) start() {
	rc.wg.Add(1)
	go func() {
		defer rc.wg.Done()
		go rc.startInput()
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

func openChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}
	return ch, nil
}
