package socketio

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// Parse raw url string and make valid handshake or websockets socket.io url
type urlParser struct {
	raw    string
	parsed *url.URL
	path   string
	token  string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newURLParser(raw string) (*urlParser, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme == "" {
		parsed.Scheme = "http"
	}

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 7)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return &urlParser{raw: raw, parsed: parsed, token: string(b)}, nil
}

func (u *urlParser) origin() string {
	return u.parsed.String()
}

func (u *urlParser) handshake() string {
	return fmt.Sprintf("%s/%s/?EIO=3&transport=polling&t=%s", u.parsed.String(), u.path, u.token)
}

func (u *urlParser) websocket(sessionId string) string {
	host := strings.Replace(u.parsed.String(), "http://", "ws://", 1)

	if u.parsed.Scheme == "https" {
		host = strings.Replace(u.parsed.String(), "https://", "wss://", 1)
	}

	return fmt.Sprintf("%s/%s/?EIO=3&transport=websocket&sid=%s", host, u.path, sessionId)
}
