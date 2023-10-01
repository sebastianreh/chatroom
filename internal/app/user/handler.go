package user

import (
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type UserHandler interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
}

type userHandler struct {
	config  config.Config
	service UserService
	logs    logger.Logger
}

func NewUserHandler(cfg config.Config, service UserService, logger logger.Logger) UserHandler {
	return &userHandler{
		config:  cfg,
		service: service,
		logs:    logger,
	}
}

func (u userHandler) Create(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHandler) Login(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHandler) Get(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userHandler) Delete(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
