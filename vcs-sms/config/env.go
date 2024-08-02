package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
	DB_PORT     string `env:"DB_PORT" envDefault:"5432"`
	DB_USER     string `env:"DB_USER" envDefault:"checkpoint-project-client"`
	DB_PASS     string `env:"DB_PASS,required"`
	DB_DATABASE string `env:"DB_DATABASE" envDefault:"checkpoint-vcs-sms-db"`

	REDIS_HOST string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	REDIS_PWD  string `env:"REDIS_PWD"`

	ELASTIC_ENDPOINT string `env:"ELASTIC_ENDPOINT" envDefault:"http://localhost:9200"`

	KAFKA_BOOTSTRAP_SERVERS        string `env:"KAFKA_BOOTSTRAP_SERVERS"`
	KAFKA_SSL_CA_LOCATION          string `env:"KAFKA_SSL_CA_LOCATION"`
	KAFKA_SSL_CERTIFICATE_LOCATION string `env:"KAFKA_SSL_CERTIFICATE_LOCATION"`
	KAFKA_SSL_KEY_LOCATION         string `env:"KAFKA_SSL_KEY_LOCATION"`
	KAFKA_SSL_KEY_PASSWORD         string `env:"KAFKA_SSL_KEY_PASSWORD"`

	GRPC_HOST string `env:"GRPC_HOST"`
	GRPC_PORT string `env:"GRPC_PORT"`
}

var once sync.Once
var cfg *Config

func InitConfig() *Config {
	if cfg == nil {
		once.Do(func() {
			cfg = &Config{}
			if err := env.Parse(cfg); err != nil {
				panic(err)
			}
		})
	}
	fmt.Printf("Config: %+v\n", cfg)
	return cfg
}

func GetEnv() Config {
	return *cfg
}
