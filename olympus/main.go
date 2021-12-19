package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	brokerConfig "github.com/project-auxo/auxo/olympus/internal/config"
	"github.com/project-auxo/auxo/olympus/logging"
	"github.com/project-auxo/auxo/olympus/pkg/broker"
	"github.com/project-auxo/auxo/olympus/pkg/util"
)

var log = logging.Base()

func readConf(configPath string) *brokerConfig.Config {
	util.Validate(configPath)
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("fail to read config path: %v", err)
	}
	cfg := &brokerConfig.Config{}
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

	broker, err := broker.New(cfg)
	if err != nil {
		log.Fatalln("Failed to start the broker")
	}
	broker.Run()
}
