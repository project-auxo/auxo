package config

type Config struct {
	Agent struct {
		Name    string `yaml:"name"`
		Olympus string `yaml:"olympus"`
		Port    int    `yaml:"port"`
	} `yaml:"agent"`
}
