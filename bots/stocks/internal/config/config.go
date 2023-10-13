package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Kafka struct {
			Server      string `envconfig:"KAFKA_SERVER" default:"chatroom"`
			StocksTopic string `envconfig:"STOCKS_TOPIC" default:"chatroom"`
		}
		Websocket struct {
			Host string `envconfig:"REDIS_HOST" default:"localhost:6379"`
		}
	}
)

var (
	Configs Config
)

func NewConfig() Config {
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
