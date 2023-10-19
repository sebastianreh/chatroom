package room

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

const handlerName = "room.handler"

type RoomHandler interface {
	Create(c echo.Context) error
	Join(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
}

type roomHandler struct {
	config  config.Config
	service RoomService
	logs    logger.Logger
}

func NewRoomHandler(cfg config.Config, service RoomService, logger logger.Logger) RoomHandler {
	return &roomHandler{
		config:  cfg,
		service: service,
		logs:    logger,
	}
}

func (handler *roomHandler) Create(ctx echo.Context) error {
	request := new(entities.Room)
	if err := ctx.Bind(request); err != nil {
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Set"))
		ctx.Error(err)
		return nil
	}

	response, err := handler.service.Create(ctx.Request().Context(), *request)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (handler *roomHandler) Get(ctx echo.Context) error {
	roomSearch := new(entities.RoomSearch)
	if err := ctx.Bind(roomSearch); err != nil {
		err = resterror.NewBadRequestError(err.Error())
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Get"))
		ctx.Error(err)
		return nil
	}

	rooms, err := handler.service.Get(ctx.Request().Context(), *roomSearch)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.JSON(http.StatusOK, rooms)
}

func (handler *roomHandler) Delete(ctx echo.Context) error {
	roomID := ctx.Param("id")
	if str.IsEmpty(roomID) {
		err := errors.New("error: empty id")
		handler.logs.Error(str.ErrorConcat(err, handlerName, "Delete"))
		ctx.Error(err)
		return nil
	}

	err := handler.service.Delete(ctx.Request().Context(), roomID)
	if err != nil {
		ctx.Error(err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (handler *roomHandler) Join(c echo.Context) error {
	//roomSearch := new(entities.RoomSearch)
	return nil
}
