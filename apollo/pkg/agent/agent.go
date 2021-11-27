package agent

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	zmq "github.com/pebbe/zmq4"

	agentCfg "github.com/project-auxo/auxo/apollo/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	discpb "github.com/project-auxo/auxo/olympus/proto/discovery"
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

	agent.log.Infof("%s is running...", agent.name)
	go agent.actor.run()
	interrupt := <-runChan
	agent.log.Infof("Agent is shutting down due to %+v\n", interrupt)
	return
}
