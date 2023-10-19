package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Kafka struct {
			Server      string `envconfig:"KAFKA_SERVER" default:"localhost:9092"`
			StocksTopic string `envconfig:"KAFKA_STOCKS_TOPIC" default:"stocks"`
			MaxRetries  int    `envconfig:"KAFKA_MAX_RETRIES" default:"3"`
		}
		Websocket struct {
			Endpoint string `envconfig:"WEBSOCKET_ENDPOINT" default:"ws://localhost:8000/chatroom/session/bot"`
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
