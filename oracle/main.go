package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"

	"github.com/project-auxo/auxo/olympus/logging"
	"github.com/project-auxo/auxo/olympus/pkg/util"
	oracle "github.com/project-auxo/auxo/oracle/grpc"
	oracleCfg "github.com/project-auxo/auxo/oracle/internal/config"
	pb "github.com/project-auxo/auxo/oracle/proto"
)

var log = logging.Base()

func readConf(configPath string) *oracleCfg.Config {
	util.Validate(configPath)
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("fail to read config path: %v", err)
	}
	cfg := &oracleCfg.Config{}
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		log.Fatalf("in file %q: %v", configPath, err)
	}
	return cfg
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()
	cfg := readConf(configPath)

	var runChan = make(chan os.Signal, 1)
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	lis, err := net.Listen("tcp", fmt.Sprintf(
		"%s:%d", cfg.Oracle.Hostname, cfg.Oracle.Port))
	if err != nil {
		log.Fatalf("set up oracle server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOracleBackendServiceServer(s, &oracle.OracleBackendServer{})

	log.Info("Oracle server is running...")
	go func() {
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-runChan
}
