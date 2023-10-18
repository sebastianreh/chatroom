package container

import (
	"errors"
	"flag"
	"github.com/go-resty/resty/v2"
	"github.com/sebastianreh/chatroom-bots/stocks/internal/config"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/kafka"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/rest"
	"github.com/sebastianreh/chatroom-bots/stocks/pkg/websocket"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type Container struct {
	Config   config.Config
	Client   rest.StooqClient
	Producer kafka.Producer
	Socket   websocket.Websocket
	Logs     logger.Logger
	RoomID   string
}

func Build(botName string) Container {
	log := logger.NewLogger()
	roomID := getRoomId(log)
	cfg := config.NewConfig()
	restyClient := resty.New()
	socket := websocket.NewWebsocket(log, cfg, botName, roomID)
	producer, err := kafka.NewProducer(log, cfg.Kafka.Server, kafka.WithMaxRetries(cfg.Kafka.MaxRetries))
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	log.Info("Succesfully started bot on roomID:", roomID)

	return Container{
		Config:   cfg,
		Logs:     log,
		Client:   rest.NewStooqClient(log, restyClient),
		Producer: producer,
		Socket:   socket,
		RoomID:   roomID,
	}
}

func getRoomId(log logger.Logger) string {
	roomID := flag.String("room_id", "", "Room ID")

	flag.Parse()

	if roomID == nil || *roomID == "" {
		err := errors.New("empty room_id")
		log.Fatal(err.Error())
		panic(err)
	}

	return *roomID
}
