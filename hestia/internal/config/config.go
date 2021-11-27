package config

type Config struct {
	Hestia struct {
		Hostname       string `yaml:"hostname"`
		Port           int    `yaml:"port"`
		FrontendClient struct {
			Hostname string `yaml:"hostname"`
			Port     int    `yaml:"port"`
		} `yaml:"frontend_client"`
	} `yaml:"hestia"`
}
