package socketio

import (
	"errors"
	"log"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func newTransport(socket *Socket) error {
	if !socket.Session.SupportProtocol("websocket") {
		return errors.New("websocket protocol is not supported by server")
	}

	c, err := websocket.NewConfig(
		socket.Session.URL.websocket(socket.Session.ID), socket.Session.URL.origin())
	if err != nil {
		return errors.New("could not create ws config: " + err.Error())
	}

	ws, err := websocket.DialConfig(c)
	if err != nil {
		return errors.New("could not dial the server: " + err.Error())
	}

	_, err = ws.Write(connectMsg().Bytes())
	if err != nil {
		return errors.New("could not write connect message: " + err.Error())
	}

	go func() {
		ticker := time.NewTicker(socket.Session.HeartbeatTimeout)
		for {
			select {
			case msg := <-socket.Send:
				ws.Write(msg.Bytes())
			case <-ticker.C:
				ws.Write(heartbeatMsg().Bytes())
			}
		}
	}()

	go func() error {
		for {
			buff := make([]byte, 16*1024)
			n, err := ws.Read(buff)
			if err != nil {
				return errors.New("error from server: " + err.Error())
			}

			body := string(buff[:n])

			if strings.HasPrefix(body, "3probe") {
				socket.Send <- ackMsg()
			}

			// This is a heartbeat reply, ignore
			if strings.HasPrefix(body, HeartBeatReply) {
				continue
			}

			if strings.HasPrefix(body, MsgFromServer) {
				msg, err := parseMessage(buff[:n])
				if err != nil {
					return errors.New("failed to parse the response from server: " + err.Error())
				}

				socket.Receive <- msg
			} else {
				log.Println("received something that dont know how to parse")
			}
		}
	}()

	return nil
}
