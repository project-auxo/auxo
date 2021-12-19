package config

type Config struct {
	Oracle struct {
		Name     string `yaml:"name"`
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
	} `yaml:"oracle"`
}
