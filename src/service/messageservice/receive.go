package messageservice

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/param/messageparam"
	"net"
	"time"
)

func (s Service) Receive(req messageparam.ReceiveRequest) error {
	chanel := s.rabbitmq.GetOutputChannel(req.UserId)
	heartBeat := s.rabbitmq.GetHeartbeatChannel(req.UserId)
	fmt.Println("Start web socket connection")

	conn, _, _, err := ws.UpgradeHTTP(req.Request, req.Response)
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	websocketClosedChannel := make(chan bool)
	go readMessage(conn, websocketClosedChannel)

	go func() {
		for {
			select {
			case msg := <-chanel:
				body := msg.Body
				err := wsutil.WriteServerMessage(conn, ws.OpText, body)
				if err != nil {
					logger.Logger.Error(err.Error())
					return
				}
				err = msg.Ack()
				if err != nil {
					logger.Logger.Error(err.Error())
					return
				}
			case <-websocketClosedChannel:
				return
			case <-time.After(time.Second):
				err := wsutil.WriteServerMessage(conn, ws.OpPing, []byte("ping"))
				if err != nil {
					logger.Logger.Error(err.Error())
					return
				}
				heartBeat <- true
			}
		}
	}()

	return nil
}

func readMessage(conn net.Conn, websocketClosedChannel chan<- bool) {
	for {
		_, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			websocketClosedChannel <- true
			return
		}
	}
}
