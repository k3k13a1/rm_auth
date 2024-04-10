package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env             string        `yaml:"env"`
	TokenTTL        time.Duration `yaml:"token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	JWTSecret       string        `yaml:"jwt_secret"`
	REST            RESTConfig    `yaml:"rest"`
}

type RESTConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func SetupConfig() *Config {
	var config Config

	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		panic("config zapili blya")
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic("config hueta, peredelai")
	}

	return &config
}
