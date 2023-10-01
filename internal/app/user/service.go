package user

import (
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

type UserService interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
}

type userService struct {
	config     config.Config
	repository UserRepository
	logs       logger.Logger
}

func NewUserService(cfg config.Config, repository UserRepository, logger logger.Logger) UserService {
	return &userService{
		config:     cfg,
		repository: repository,
		logs:       logger,
	}
}

func (u userService) Create(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) Login(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) Get(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) Delete(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
