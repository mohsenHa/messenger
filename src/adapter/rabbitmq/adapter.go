package rabbitmq

import (
	"fmt"
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
	fmt.Printf("Main rabbitmq object address %p \n", &rabbitmq)
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
		failOnError(err, "rabbitmq connection failed")
		if err == nil {
			break
		}
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
	failOnError(err, "Failed to connect to rabbitmq server")
	if err != nil {
		return err
	}
	ca.rabbitmq.connection = conn
	fmt.Println("Connected to rabbitmq server")
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
				failOnError(e, "Connection failed to rabbitmq")
				if e == nil {
					break
				}
			}

			return
		}
	}
}

func (ca *ChannelAdapter) NewChannel(name string) {
	closeSignalChannel := make(chan bool, 10)
	heartBeatSignalChannel := make(chan bool, 10)
	ca.channels[name] = newChannel(
		ca.done,
		ca.wg,
		rabbitmqChannelParams{
			rabbitmq:               ca.rabbitmq,
			exchange:               name + "-exchange",
			queue:                  name + "-queue",
			bufferSize:             ca.config.BufferSize,
			maxRetryPolicy:         ca.config.MaxRetryPolicy,
			closeSignalChannel:     closeSignalChannel,
			heartBeatSignalChannel: heartBeatSignalChannel,
		})
	ca.CloseIdleChannel(name, closeSignalChannel, heartBeatSignalChannel)
}

func (ca *ChannelAdapter) CloseIdleChannel(name string, closeSignalChannel chan<- bool, heartBeatSignalChannel <-chan bool) {
	ca.wg.Add(1)
	go func() {
		defer ca.wg.Done()
		for {
			timer := time.After(time.Second * time.Duration(ca.config.ChannelCleanerTimerInSecond))
			select {
			case <-timer:
				fmt.Printf("Channel is idle more than %d seconds\n", ca.config.ChannelCleanerTimerInSecond)
				close(closeSignalChannel)
				delete(ca.channels, name)
				return
			case <-heartBeatSignalChannel:
				fmt.Println("Heartbeat signal received")
				break
			case <-ca.done:
				return
			}
		}
	}()
}

func (ca *ChannelAdapter) GetInputChannel(name string) chan<- []byte {
	if c, ok := ca.channels[name]; ok {
		return c.GetInputChannel()
	}
	ca.NewChannel(name)
	return ca.GetInputChannel(name)
}

func (ca *ChannelAdapter) GetOutputChannel(name string) <-chan Message {
	if c, ok := ca.channels[name]; ok {
		return c.GetOutputChannel()
	}
	ca.NewChannel(name)
	return ca.GetOutputChannel(name)
}

func WaitForConnection(rabbitmq *Rabbitmq) {
	fmt.Printf("the address in wait for connection %p \n", rabbitmq)

	rabbitmq.cond.L.Lock()
	defer rabbitmq.cond.L.Unlock()
	for rabbitmq.connection.IsClosed() {
		fmt.Println(rabbitmq.connection.IsClosed())
		fmt.Println("Before wait for connection")
		rabbitmq.cond.Wait()
		fmt.Println("After wait for connection")

	}
}
