package agent

import (
	"fmt"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"

	"github.com/project-auxo/auxo/olympus/logging"
	util "github.com/project-auxo/auxo/olympus/pkg/util"
	discpb "github.com/project-auxo/auxo/olympus/proto/discovery"
)

const (
	heartbeatInterval = time.Duration(1) * time.Second
	workersEndpoint   = "inproc://workers"
)

var _readyMsg = &discpb.DiscoveryMessage{
	Header:  discpb.Header_HEADER_READY,
	Origin:  &discpb.Entity{Type: agentEntityType},
	Command: &discpb.DiscoveryMessage_Ready{},
}

type Actor struct {
	log              logging.Logger
	name             string
	externalSocket   *zmq.Socket
	externalEndpoint string
	workersSocket    *zmq.Socket // Communicate with internal workers.
	poller           *zmq.Poller
}

func newActor(name string, externalEndpoint string) (actor *Actor, err error) {
	var externalSocketErr error
	var workersSocketErr error

	actor = &Actor{log: logging.Base(), name: name, externalEndpoint: externalEndpoint, poller: zmq.NewPoller()}
	actor.externalSocket, externalSocketErr = zmq.NewSocket(zmq.DEALER)
	actor.workersSocket, workersSocketErr = zmq.NewSocket(zmq.DEALER)
	err = multierror.Append(err, externalSocketErr, workersSocketErr)
	actor.poller.Add(actor.externalSocket, zmq.POLLIN)
	actor.poller.Add(actor.workersSocket, zmq.POLLIN)
	return
}

func (actor *Actor) bind() (err error) {
	if socketErr := actor.externalSocket.Connect(actor.externalEndpoint); socketErr != nil {
		err = multierror.Append(err, socketErr)
	} else {
		actor.log.Debugf("%s sending ready message", actor.name)
		actor.send(actor.externalSocket, _readyMsg)
	}
	if socketErr := actor.workersSocket.Bind(workersEndpoint); socketErr != nil {
		err = multierror.Append(err, socketErr)
	}
	return
}

func (actor *Actor) close() (err error) {
	if actor.externalSocket != nil {
		err = multierror.Append(err, actor.externalSocket.Close())
		actor.externalSocket = nil
	}
	if actor.workersSocket != nil {
		err = multierror.Append(err, actor.workersSocket.Close())
		actor.workersSocket = nil
	}
	return
}

func (actor *Actor) run() (err error) {
	err = actor.bind()
	if err != nil {
		return
	}
	for {
		polled, err := actor.poller.Poll(heartbeatInterval)
		if err != nil {
			// Interrupted
			break
		}
		if len(polled) > 0 {
			for _, socket := range polled {
				switch s := socket.Socket; s {
				case actor.externalSocket:
					actor.handleExternalSocket()
				case actor.workersSocket:
					actor.handleWorkersSocket()
				}
			}
		}
	}
	return
}

func (actor *Actor) handleExternalSocket() (err error) {
	recvBytes, recvErr := actor.externalSocket.RecvBytes(0)
	if recvErr != nil {
		err = multierror.Append(err, recvErr)
	}
	msg, unmarshalErr := util.UnmarshalDiscoveryMessage(recvBytes)
	if unmarshalErr != nil {
		err = multierror.Append(err, unmarshalErr)
	}
	fmt.Println(msg)
	return
}

func (actor *Actor) handleWorkersSocket() {
	// Empty
}

func (actor *Actor) send(socket *zmq.Socket, msg *discpb.DiscoveryMessage) (err error) {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	_, err = socket.SendBytes(msgBytes, zmq.DONTWAIT)
	return
}
