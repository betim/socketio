package socketio

import (
	"bytes"
	"fmt"
)

const HeartBeatReply = "3"
const MsgFromServer = "42["

type Message struct {
	Type     int
	Endpoint string
	Data     string
}

func parseMessage(body []byte) (*Message, error) {
	// body[bytes.Index(body, []byte("[")):]
	msg := Message{
		Type:     42,
		Endpoint: string(body[bytes.Index(body, []byte("["))+2 : bytes.Index(body, []byte(","))-1]),
		Data:     string(body[bytes.Index(body, []byte(","))+2 : bytes.Index(body, []byte("]"))]),
	}

	return &msg, nil
}

// String returns the string represenation of the Message
func (m Message) String() string {
	return fmt.Sprintf("%d%s", m.Type, m.Data)
}

// Bytes is a handy func to convert string to []byte for ws.socket
func (m Message) Bytes() []byte {
	return []byte(m.String())
}

func connectMsg() *Message {
	return &Message{Type: 2, Data: "probe"}
}

func heartbeatMsg() *Message {
	return &Message{Type: 2}
}

// NewMessage given an endpoint and message will construct a new message for you
func NewMessage(endpoing, msg string) *Message {
	return &Message{Type: 42, Data: msg}
}

func ackMsg() *Message {
	return &Message{Type: 5}
}
