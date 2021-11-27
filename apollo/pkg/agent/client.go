package agent

import (
	"errors"

	zmq "github.com/pebbe/zmq4"
)

var errPermanent = errors.New("permanent error, abandoning request")

type Client struct {
	socket  *zmq.Socket
	olympus string // Where to connect to Olympus
}
