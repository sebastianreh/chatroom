package container

import (
	"github.com/sebastianreh/chatroom/internal/app/ping"
	"github.com/sebastianreh/chatroom/internal/app/room"
	"github.com/sebastianreh/chatroom/internal/app/session"
	"github.com/sebastianreh/chatroom/internal/app/user"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/kafka"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/pkg/mongodb"
	rds "github.com/sebastianreh/chatroom/pkg/redis"
	ws "github.com/sebastianreh/chatroom/pkg/websocket"
)

type Dependencies struct {
	PingHandler    ping.Handler
	Config         config.Config
	Logs           logger.Logger
	UserHandler    user.UserHandler
	RoomHandler    room.RoomHandler
	SessionHandler session.SessionHandler
}

func Build() Dependencies {
	dependencies := Dependencies{}
	dependencies.Config = config.NewConfig()
	logs := logger.NewLogger()
	dependencies.Logs = logs
	dependencies.PingHandler = ping.NewSHandierPing(dependencies.Config)

	mongoDB := mongodb.NewMongoDB(dependencies.Config)
	redis, err := rds.NewRedis(logs, dependencies.Config)
	if err != nil {
		logs.Fatal(err.Error())
	}
	websocket := ws.NewWebsocket()
	kafkaConsumer, err := kafka.NewKafkaConsumer(dependencies.Config, dependencies.Logs)
	if err != nil {
		logs.Fatal(err.Error())
	}

	userRepository := user.NewUserRepository(dependencies.Config, mongoDB, dependencies.Logs)
	userService := user.NewUserService(dependencies.Config, userRepository, dependencies.Logs)
	userHandler := user.NewUserHandler(dependencies.Config, userService, dependencies.Logs)

	roomRepository := room.NewRoomRepository(dependencies.Config, mongoDB, dependencies.Logs)
	roomService := room.NewRoomService(dependencies.Config, roomRepository, dependencies.Logs)
	roomHandler := room.NewRoomHandler(dependencies.Config, roomService, dependencies.Logs)

	sessionRepository := session.NewSessionRepository(dependencies.Config, redis, dependencies.Logs)
	sessionService := session.NewSessionService(dependencies.Config, sessionRepository, dependencies.Logs)
	sessionHandler := session.NewSessionHandler(dependencies.Config, sessionService, websocket, dependencies.Logs)

	dependencies.UserHandler = userHandler
	dependencies.RoomHandler = roomHandler
	dependencies.SessionHandler = sessionHandler

	return dependencies
}
