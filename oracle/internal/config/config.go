package config

type Config struct {
	Oracle struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
	} `yaml:"oracle"`
}
