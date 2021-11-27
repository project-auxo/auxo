package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	agentCfg "github.com/project-auxo/auxo/apollo/internal/config"
	agent "github.com/project-auxo/auxo/apollo/pkg/agent"
	"github.com/project-auxo/auxo/olympus/logging"
	"github.com/project-auxo/auxo/olympus/pkg/util"
)

var log = logging.Base()

func readConf(configPath string) *agentCfg.Config {
	util.Validate(configPath)
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("fail to read config path: %v", err)
	}
	cfg := &agentCfg.Config{}
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

	agent := agent.New(cfg)
	agent.Run()
}
