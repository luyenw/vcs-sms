package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	KAFKA_BOOTSTRAP_SERVERS        string `env:"KAFKA_BOOTSTRAP_SERVERS"`
	KAFKA_SSL_CA_LOCATION          string `env:"KAFKA_SSL_CA_LOCATION"`
	KAFKA_SSL_CERTIFICATE_LOCATION string `env:"KAFKA_SSL_CERTIFICATE_LOCATION"`
	KAFKA_SSL_KEY_LOCATION         string `env:"KAFKA_SSL_KEY_LOCATION"`
	KAFKA_SSL_KEY_PASSWORD         string `env:"KAFKA_SSL_KEY_PASSWORD"`

	MAIL_PASSWORD string `env:"MAIL_PASSWORD"`
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
