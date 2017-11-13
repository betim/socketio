package socketio

import "errors"

// Socket holds info with whom and as who is connected and 2 communication chans
type Socket struct {
	URL     string
	Session *Session
	Receive chan *Message
	Send    chan *Message
}

// Stream will try to handshake and establish a new session
// then will try to establish a websocket connection and return a *Socket
func Stream(url string, path string, query string) (*Socket, error) {
	session, err := NewSession(url, path, query)
	if err != nil {
		return nil, errors.New("could not create a new session: " + err.Error())
	}

	rch := make(chan *Message)
	sch := make(chan *Message)
	s := &Socket{url, session, rch, sch}
	if err := newTransport(s); err != nil {
		return nil, err
	}

	return s, nil
}
