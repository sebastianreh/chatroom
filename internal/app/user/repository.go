package user

import (
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/pkg/mongodb"
)

type UserRepository interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
}

type userRepository struct {
	config  config.Config
	mongodb mongodb.MongoDBier
	logs    logger.Logger
}

func NewUserRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logger.Logger) UserRepository {
	return &userRepository{
		config:  cfg,
		mongodb: mongoDBier,
		logs:    logger,
	}
}

func (u userRepository) Create(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) Login(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) Get(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepository) Delete(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
