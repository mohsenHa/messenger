package rabbitmq

import (
	"fmt"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type ChannelAdapter struct {
	wg                       *sync.WaitGroup
	done                     <-chan bool
	channels                 map[string]*rabbitmqChannel
	config                   Config
	rabbitmq                 *Rabbitmq
	rabbitmqConnectionClosed chan bool
}
type Rabbitmq struct {
	connection *amqp.Connection
	cond       *sync.Cond
}

func New(done <-chan bool, wg *sync.WaitGroup, config Config) *ChannelAdapter {
	cond := sync.NewCond(&sync.Mutex{})
	rabbitmq := Rabbitmq{
		cond:       cond,
		connection: &amqp.Connection{},
	}
	c := &ChannelAdapter{
		done:                     done,
		wg:                       wg,
		config:                   config,
		rabbitmq:                 &rabbitmq,
		channels:                 make(map[string]*rabbitmqChannel),
		rabbitmqConnectionClosed: make(chan bool),
	}

	for {
		err := c.connect()
		time.Sleep(time.Second * time.Duration(config.ReconnectSecond))
		if err == nil {
			break
		}
		logger.NewLog("rabbitmq connection error").
			WithCategory(loggerentity.CategoryRabbitMQ).
			WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

	}

	return c
}

func (ca *ChannelAdapter) connect() error {
	ca.rabbitmq.cond.L.Lock()
	defer ca.rabbitmq.cond.L.Unlock()
	close(ca.rabbitmqConnectionClosed)
	ca.rabbitmqConnectionClosed = make(chan bool)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		ca.config.User, ca.config.Password, ca.config.Host,
		ca.config.Port, ca.config.Vhost))
	if err != nil {
		return err
	}
	ca.rabbitmq.connection = conn
	logger.NewLog("Connected to rabbitmq server").
		WithCategory(loggerentity.CategoryRabbitMQ).
		WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
		Debug()

	ca.rabbitmq.cond.Broadcast()

	ca.wg.Add(1)
	go func() {
		defer ca.wg.Done()
		for {
			select {
			case <-ca.done:

				return
			case <-ca.rabbitmqConnectionClosed:

				return
			}
		}
	}()
	go ca.waitForConnectionClose()

	return nil
}

func (ca *ChannelAdapter) waitForConnectionClose() {
	connectionClosedChannel := make(chan *amqp.Error)
	ca.rabbitmq.connection.NotifyClose(connectionClosedChannel)

	for {
		select {
		case <-ca.done:
			return
		case err := <-connectionClosedChannel:
			fmt.Println("Connection closed")
			fmt.Println(err)
			for {
				e := ca.connect()
				time.Sleep(time.Second * time.Duration(ca.config.ReconnectSecond))
				if e == nil {
					break
				}
				logger.NewLog("Connection failed to rabbitmq").
					WithCategory(loggerentity.CategoryRabbitMQ).
					WithSubCategory(loggerentity.SubCategoryRabbitMQConnection).
					With(loggerentity.ExtraKeyErrorMessage, err.Error()).
					Error()
			}

			return
		}
	}
}

func (ca *ChannelAdapter) NewChannel(name string) error {
	chanBuffer := 10
	rabbitMQCloseSignal := make(chan bool, chanBuffer)
	noOutputConsumer := make(chan bool, chanBuffer)
	heartBeatSignalChannel := make(chan bool, chanBuffer)
	channel, err := newChannel(
		ca.done,
		ca.wg,
		rabbitmqChannelParams{
			rabbitmq:               ca.rabbitmq,
			exchange:               name + "-exchange",
			queue:                  name + "-queue",
			bufferSize:             ca.config.BufferSize,
			maxRetryPolicy:         ca.config.MaxRetryPolicy,
			rabbitMQCloseSignal:    rabbitMQCloseSignal,
			noOutputConsumer:       noOutputConsumer,
			heartBeatSignalChannel: heartBeatSignalChannel,
		})
	if err != nil {
		return err
	}
	ca.channels[name] = channel
	ca.CloseIdleChannel(name, rabbitMQCloseSignal, noOutputConsumer, heartBeatSignalChannel)

	return nil
}

func (ca *ChannelAdapter) CloseIdleChannel(name string, closeSignalChannel chan bool, noOutputConsumer chan bool, heartBeatSignalChannel <-chan bool) {
	ca.wg.Add(1)
	go func() {
		defer ca.wg.Done()
		for {
			timer := time.After(time.Second * time.Duration(ca.config.ChannelCleanerTimerInSecond))
			select {
			case <-timer:
				logger.L().Debugf("Channel is idle more than %d seconds. produce close signal.", ca.config.ChannelCleanerTimerInSecond)
				close(closeSignalChannel)
				delete(ca.channels, name)

				return
			case <-noOutputConsumer:
				logger.L().Debugf("Channel there is no output consumer produce close signal")
				close(closeSignalChannel)
				delete(ca.channels, name)

				return
			case <-heartBeatSignalChannel:
				break
			case <-closeSignalChannel:
				delete(ca.channels, name)

				return
			case <-ca.done:
				return
			}
		}
	}()
}

func (ca *ChannelAdapter) GetInputChannel(name string) (chan<- []byte, error) {
	if c, ok := ca.channels[name]; ok {
		return c.GetInputChannel(), nil
	}
	err := ca.NewChannel(name)
	if err != nil {
		return nil, err
	}

	return ca.GetInputChannel(name)
}

func (ca *ChannelAdapter) GetOutputChannel(name string, outputChannelCloseSignal <-chan bool) (<-chan Message, error) {
	if c, ok := ca.channels[name]; ok {
		return c.GetOutputChannel(outputChannelCloseSignal), nil
	}
	err := ca.NewChannel(name)
	if err != nil {
		return nil, err
	}

	return ca.GetOutputChannel(name, outputChannelCloseSignal)
}

func (ca *ChannelAdapter) GetHeartbeatChannel(name string) (chan<- bool, error) {
	if c, ok := ca.channels[name]; ok {
		return c.GetHeartbeatChannel(), nil
	}
	err := ca.NewChannel(name)
	if err != nil {
		return nil, err
	}

	return ca.GetHeartbeatChannel(name)
}

func WaitForConnection(rabbitmq *Rabbitmq) {
	rabbitmq.cond.L.Lock()
	defer rabbitmq.cond.L.Unlock()
	for rabbitmq.connection.IsClosed() {
		rabbitmq.cond.Wait()
	}
}
