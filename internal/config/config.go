package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv          string        `mapstructure:"APP_ENV"`
	TokenTTL        time.Duration `mapstructure:"TOKEN_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`
	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	RESTHost        string        `mapstructure:"REST_HOST"`
	RESTPort        int           `mapstructure:"REST_PORT"`
	RESTTimeout     time.Duration `mapstructure:"REST_TIMEOUT"`
	GRPCPort        int           `mapstructure:"GRPC_PORT"`
	PSQLHost        string        `mapstructure:"PSQL_HOST"`
	PSQLPort        int           `mapstructure:"PSQL_PORT"`
	PSQLUser        string        `mapstructure:"PSQL_USER"`
	PSQLPass        string        `mapstructure:"PSQL_PASS"`
	PSQLDB          string        `mapstructure:"PSQL_DB"`
	PSQLSSLMode     string        `mapstructure:"PSQL_SSLMODE"`
	PSQLTimeout     time.Duration `mapstructure:"PSQL_TIMEOUT"`
}

func SetupConfig() *Config {
	cfg := Config{}

	viper.AutomaticEnv()
	viper.SetConfigFile("./config/.env")

	err := viper.ReadInConfig()
	if err != nil {
		panic("did not find config file")
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.AppEnv == "development" {
		fmt.Println("running in development mode")
	}

	return &cfg
}
