package user

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/cmd/httpserver/resterror"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"net/http"
)

const handlerName = "user.handler"

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

func (handler *userHandler) Create(ctx echo.Context) error {
	request := new(entities.User)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Create"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Create(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusCreated)
}

func (handler *userHandler) Login(ctx echo.Context) error {
	request := new(entities.User)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Login"))
		ctx.Error(err)
		return nil
	}

	response, err := handler.service.Login(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, response)
}

func (handler *userHandler) Get(ctx echo.Context) error {
	userSearch := new(entities.UserSearch)
	if err := ctx.Bind(userSearch); err != nil {
		err = resterror.NewBadRequestError(err.Error())
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Get"))
		ctx.Error(err)
		return nil
	}

	users, err := handler.service.Get(ctx.Request().Context(), *userSearch)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, users)
}

func (handler *userHandler) Delete(ctx echo.Context) error {
	userID := ctx.Param("id")
	if str.IsEmpty(userID) {
		err := errors.New("error: empty id")
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Delete(ctx.Request().Context(), userID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}
