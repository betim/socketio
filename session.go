package socketio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Session holds the configuration variables received from the server
type Session struct {
	ID                 string        `json:"sid"`
	HeartbeatTimeout   time.Duration `json:"pingInterval"`
	ConnectionTimeout  time.Duration `json:"pingTimeout"`
	SupportedProtocols []string      `json:"upgrades"`
	URL                *urlParser
}

// NewSession receives the configuraiton variables from the server
func NewSession(url, path, query string) (*Session, error) {
	urlParser, err := newURLParser(url)
	if err != nil {
		return nil, err
	}

	urlParser.path = path
	response, err := http.Get(urlParser.handshake())
	if err != nil {
		return nil, errors.New("http.Get: " + err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("could not read body: " + err.Error())
	}
	response.Body.Close()

	s := Session{
		URL: urlParser,
	}

	// for some reason, socket.io sends some garbage at the begining of the response
	err = json.Unmarshal(body[bytes.Index(body, []byte("{")):], &s)
	if err != nil {
		return nil, errors.New("failed to handshake with the server: " + err.Error())
	}

	s.HeartbeatTimeout *= time.Second
	s.ConnectionTimeout *= time.Second

	return &s, nil
}

// SupportProtocol checks if the given protocol is supported by the server
func (session *Session) SupportProtocol(protocol string) bool {
	for _, supportedProtocol := range session.SupportedProtocols {
		if protocol == supportedProtocol {
			return true
		}
	}

	return false
}
