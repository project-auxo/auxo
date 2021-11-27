package config

type Config struct {
	Broker struct {
		Hostname       string `yaml:"hostname"`
		Port           int    `yaml:"port"`
		FrontendServer struct {
			Hostname string `yaml:"hostname"`
			Port     int    `yaml:"port"`
		} `yaml:"frontend_server"`
	} `yaml:"broker"`
}
