package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"

	hestiaCfg "github.com/project-auxo/auxo/hestia/internal/config"
	authenticator "github.com/project-auxo/auxo/hestia/platform/auth"
	"github.com/project-auxo/auxo/hestia/platform/router"
	"github.com/project-auxo/auxo/olympus/logging"
	"github.com/project-auxo/auxo/olympus/pkg/util"
)

var log = logging.Base()

func readConf(configPath string) *hestiaCfg.Config {
	util.Validate(configPath)
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("fail to read config path: %v", err)
	}
	cfg := &hestiaCfg.Config{}
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		log.Fatalf("in file %q: %v", configPath, err)
	}
	return cfg
}

func main() {
	envPathPtr := flag.String("env", "./.env", "path to the .env file")
	var configPath string
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")
	flag.Parse()
	cfg := readConf(configPath)
	if err := godotenv.Load(*envPathPtr); err != nil {
		log.Fatalf("load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("initialize the authenticator: %v", err)
	}
	r := router.New(auth, cfg)
	log.Infoln("Starting Hestia")
	r.Run(fmt.Sprintf("%s:%d", cfg.Hestia.Hostname, cfg.Hestia.Port))
}
