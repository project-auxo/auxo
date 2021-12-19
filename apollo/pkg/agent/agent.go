package agent

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/protobuf/proto"

	agentCfg "github.com/project-auxo/auxo/apollo/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	discpb "github.com/project-auxo/auxo/olympus/proto/discovery"
	"github.com/project-auxo/auxo/oracle/services/auxo/seek"
	seekpb "github.com/project-auxo/auxo/oracle/services/auxo/seek/proto"
)

const (
	agentEntityType = discpb.Entity_AGENT
)

type ZeroMqSender interface {
	send(socket *zmq.Socket, msg *discpb.DiscoveryMessage) (err error)
}

type Agent struct {
	log     logging.Logger
	name    string
	olympus string // Where to connect to Olympus
	actor   *Actor
}

func New(cfg *agentCfg.Config) (agent *Agent) {
	olympus := fmt.Sprintf("tcp://%s:%d", cfg.Agent.Olympus, cfg.Agent.Port)
	agent = &Agent{log: logging.Base(), name: cfg.Agent.Name, olympus: olympus}
	agent.actor, _ = newActor(agent.name, olympus)
	return
}

func (agent *Agent) close() {
	agent.actor.close()
}

func (agent *Agent) Run() (err error) {
	var runChan = make(chan os.Signal, 1)
	defer agent.close()
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	agent.log.Infof("⇨ Auxo agent %s is running\n", agent.name)
	go agent.actor.run()
	// TODO(bellabah): Testing, not to be shipped to production
	go agent.playSeekGame()

	interrupt := <-runChan
	agent.log.Infof(
		"Auxo agent %s is shutting down due to %+v\n", agent.name, interrupt)
	return
}

// TODO(bellabah): Testing, not to be shipped to production
func (agent *Agent) playSeekGame() {
	hostname := seek.Hostname
	if hostname == "*" {
		hostname = "localhost"
	}

	stateSubscriberSocket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		agent.log.Fatalln("state receiver broken")
	}
	defer stateSubscriberSocket.Close()
	stateSubscriberSocket.SetSubscribe(seek.StateTopic)
	stateSubscriberSocket.Connect(fmt.Sprintf("tcp://%s:%d", hostname, seek.Port))

	commandSocket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		agent.log.Fatalln("command issuer broken")
	}
	defer commandSocket.Close()
	commandSocket.Connect(fmt.Sprintf("tcp://%s:%d", hostname, seek.CommandPort))
	command := &seekpb.Command{}

	fps := time.NewTicker(seek.PublishRate)
	for {
		stateMsg, _ := stateSubscriberSocket.RecvBytes(zmq.DONTWAIT)
		if len(stateMsg) > 0 {
			// Received the sim state
			state := &seekpb.SimState{}
			if err := proto.Unmarshal(stateMsg, state); err != nil {
				continue
			}

			// Send command based on information from the sim state
			if state.GoalPos.X < state.Cart.CartPos.X {
				command.Direction = seekpb.Direction_LEFT
			} else {
				command.Direction = seekpb.Direction_RIGHT
			}
			commandBytes, _ := proto.Marshal(command)
			commandSocket.SendBytes(commandBytes, zmq.DONTWAIT)
			// Just empty receive, don't care about game server's reply
			commandSocket.Recv(zmq.DONTWAIT)
		}
		<-fps.C
	}
}
