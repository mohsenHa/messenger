package messageservice

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mohsenHa/messenger/adapter/rabbitmq"
	"github.com/mohsenHa/messenger/logger"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"github.com/mohsenHa/messenger/param/messageparam"
	"io"
	"time"
)

func (s Service) Receive(req messageparam.ReceiveRequest) error {
	outputChannelCloseSignalSize := 10
	outputChannelCloseSignal := make(chan bool, outputChannelCloseSignalSize)
	chanel, err := s.rabbitmq.GetOutputChannel(req.UserID, outputChannelCloseSignal)
	if err != nil {
		return err
	}
	heartBeat, err := s.rabbitmq.GetHeartbeatChannel(req.UserID)
	if err != nil {
		return err
	}

	logger.L().Debugf("Start web socket connection")

	conn, _, _, err := ws.UpgradeHTTP(req.Request, req.Response)
	if err != nil {
		return err
	}

	websocketClosedChannel := make(chan bool)
	go readMessage(conn, websocketClosedChannel)

	go func() {
		for {
			select {
			case msg := <-chanel:
				s.processMessage(msg, conn)
			case <-websocketClosedChannel:
				go s.handleMessagesAlreadyReceivedFromRabbit(req.UserID, chanel)
				close(outputChannelCloseSignal)
				logger.L().Debugf("websocket connection closed")

				return
			case <-time.After(time.Second):
				err := wsutil.WriteServerMessage(conn, ws.OpPing, []byte("ping"))
				if err != nil {
					logger.NewLog("Error on send ping to websocket").
						With(loggerentity.ExtraKeyErrorMessage, err.Error()).
						Error()

					return
				}
				heartBeat <- true
			}
		}
	}()

	return nil
}

func (s Service) handleMessagesAlreadyReceivedFromRabbit(userID string, chanel <-chan rabbitmq.Message) {
	inputChannel, err := s.rabbitmq.GetInputChannel(userID)
	if err != nil {
		logger.NewLog("Error on get input channel").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()
	}
	for msg := range chanel {
		inputChannel <- msg.Body
		err = msg.Ack()
		if err != nil {
			logger.NewLog("Error on ACK").
				With(loggerentity.ExtraKeyErrorMessage, err.Error()).
				Error()

			continue
		}
	}
}

func (s Service) processMessage(msg rabbitmq.Message, conn io.Writer) {
	body := msg.Body
	err := wsutil.WriteServerMessage(conn, ws.OpText, body)
	if err != nil {
		logger.NewLog("Error on write to websocket").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

		return
	}

	sendMessage := messageparam.SendMessage{}
	err = json.Unmarshal(body, &sendMessage)
	if err != nil {
		logger.NewLog("Error on unmarshal message").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

		return
	}
	if sendMessage.Type == messageparam.SendMessageMessageType {
		go s.sendDelivery(sendMessage)
	}

	err = msg.Ack()
	if err != nil {
		logger.NewLog("Error on ACK").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

		return
	}
}

func (s Service) sendDelivery(sendMessage messageparam.SendMessage) {
	deliverMessage := messageparam.SendMessage{
		ID:   sendMessage.ID,
		Type: messageparam.SendMessageDeliverType,
		From: messageparam.SendUser{
			ID:        sendMessage.To.ID,
			PublicKey: sendMessage.To.PublicKey,
		},
		To: messageparam.SendUser{
			ID:        sendMessage.From.ID,
			PublicKey: sendMessage.From.PublicKey,
		},
		SendTime: time.Now().UTC(),
	}
	messageByte, err := json.Marshal(deliverMessage)
	if err != nil {
		logger.NewLog("Error on marshal to json").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

		return
	}
	fromChanel, err := s.rabbitmq.GetInputChannel(deliverMessage.To.ID)
	if err != nil {
		logger.NewLog("Error on get Input channel").
			With(loggerentity.ExtraKeyErrorMessage, err.Error()).
			Error()

		return
	}
	fromChanel <- messageByte
}

func readMessage(conn io.ReadWriter, websocketClosedChannel chan<- bool) {
	for {
		_, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			websocketClosedChannel <- true

			return
		}
	}
}
