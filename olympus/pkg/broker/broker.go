package broker

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	brokerConfig "github.com/project-auxo/auxo/olympus/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	util "github.com/project-auxo/auxo/olympus/pkg/util"
	discpb "github.com/project-auxo/auxo/olympus/proto/discovery"
	hestiapb "github.com/project-auxo/auxo/olympus/proto/hestia"
)

const (
	heartbeatInterval = time.Duration(1) * time.Second
	entityType        = discpb.Entity_BROKER
)

type Broker struct {
	log              logging.Logger
	socket           *zmq.Socket
	poller           *zmq.Poller
	hostname         string
	port             int
	frontendHostname string
	frontendPort     int
	endpoint         string
	entityType       discpb.Entity_Type
}

func New(cfg *brokerConfig.Config) (broker *Broker, err error) {
	hostname := cfg.Broker.Hostname
	if hostname == "localhost" {
		// Required by ZMQ
		hostname = "*"
	}
	port := cfg.Broker.Port
	endpoint := fmt.Sprintf("tcp://%s:%d", hostname, port)
	broker = &Broker{
		log:              logging.Base(),
		hostname:         hostname,
		port:             port,
		frontendHostname: cfg.Broker.FrontendServer.Hostname,
		frontendPort:     cfg.Broker.FrontendServer.Port,
		endpoint:         endpoint,
		entityType:       entityType,
		poller:           zmq.NewPoller(),
	}
	broker.socket, err = zmq.NewSocket(zmq.ROUTER)
	broker.poller.Add(broker.socket, zmq.POLLIN)
	return
}

// Bind will bind the broker instance to the given endpoint. Bind can be called
// multiple times.
func (broker *Broker) bind(endpoint string) (err error) {
	err = broker.socket.Bind(endpoint)
	if err != nil {
		broker.log.Fatalln("Failed to bind the broker")
	}
	broker.endpoint = endpoint
	return
}

// Close will cleanly close the broker's socket.
func (broker *Broker) close() (err error) {
	if broker.socket != nil {
		err = broker.socket.Close()
		broker.socket = nil
	}
	return
}

func (broker *Broker) handle() {
	for {
		polled, err := broker.poller.Poll(heartbeatInterval)
		if err != nil {
			// Interrupted
			break
		}
		if len(polled) > 0 {
			recvBytes, err := broker.socket.RecvBytes(0)
			if err != nil {
				break
			}
			discoveryMsg, err := util.UnmarshalDiscoveryMessage(recvBytes)
			_ = discoveryMsg
			broker.log.Debugln("Received a message")
		}
	}
}

func (broker *Broker) runFrontendServer() {
	var runChan = make(chan os.Signal, 1)
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	lis, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%d", broker.frontendHostname, broker.frontendPort))
	if err != nil {
		broker.log.Fatalf("set up frontend server: %v", err)
	}
	s := grpc.NewServer()
	hestiapb.RegisterHestiaFrontendServiceServer(
		s, &hestiaFrontendServer{broker: broker})

	broker.log.Info("Broker's frontend server is running...")
	go func() {
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			broker.log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-runChan
}

func (broker *Broker) Run() (err error) {
	if broker.endpoint == "" {
		return errors.New("must provide an endpoint")
	}
	broker.bind(broker.endpoint)

	var runChan = make(chan os.Signal, 1)
	defer broker.close()
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	broker.log.Infoln("Broker is running...")
	go broker.handle()
	broker.runFrontendServer()

	interrupt := <-runChan
	broker.log.Infof("Broker is shutting down due to %+v\n", interrupt)
	return
}
