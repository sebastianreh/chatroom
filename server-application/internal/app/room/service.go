package room

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/internal/entities/exceptions"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
)

const (
	serviceName = "room.service"
)

type RoomService interface {
	Create(ctx context.Context, room entities.Room) (entities.RoomCreateResponse, error)
	Get(ctx context.Context, search entities.RoomSearch) (entities.RoomsGetResponse, error)
	Delete(ctx context.Context, roomID string) error
}

type roomService struct {
	config     config.Config
	repository RoomRepository
	logs       logger.Logger
}

func NewRoomService(cfg config.Config, repository RoomRepository, logger logger.Logger) RoomService {
	return &roomService{
		config:     cfg,
		repository: repository,
		logs:       logger,
	}
}

func (service *roomService) Create(ctx context.Context, room entities.Room) (entities.RoomCreateResponse, error) {
	var roomCreateResponse entities.RoomCreateResponse
	rooms, err := service.repository.Get(ctx, entities.RoomSearch{Name: room.Name})
	if err != nil {
		return roomCreateResponse, err
	}

	if len(rooms) > 0 {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("room '%s' already exist", room.Name))
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return roomCreateResponse, err
	}

	roomID, err := service.repository.Create(ctx, room)
	if err != nil {
		return roomCreateResponse, err
	}

	roomCreateResponse.ID = roomID

	return roomCreateResponse, nil
}

func (service *roomService) Get(ctx context.Context, search entities.RoomSearch) (entities.RoomsGetResponse, error) {
	var roomsResponse entities.RoomsGetResponse
	rooms, err := service.repository.Get(ctx, search)
	if err != nil {
		return roomsResponse, err
	}

	if len(rooms) == 0 {
		err = exceptions.NewNotFoundException("rooms by filter not found")
		service.logs.Warn(str.ErrorConcat(err, repositoryName, "Get"))
		return roomsResponse, err
	}

	return entities.RoomsGetResponse{Rooms: rooms}, nil
}

func (service *roomService) Delete(ctx context.Context, roomID string) error {
	rooms, err := service.repository.Get(ctx, entities.RoomSearch{ID: roomID})
	if err != nil {
		return err
	}

	if len(rooms) == 0 {
		err = exceptions.NewNotFoundException(fmt.Sprintf("no room was found with id: %s", roomID))
		service.logs.Warn(str.ErrorConcat(err, serviceName, "Delete"))
		return err
	}

	if len(rooms) != 1 {
		err = errors.New("found more than one room to delete")
		service.logs.Error(str.ErrorConcat(err, serviceName, "Delete"))
		return err
	}

	rooms[0].IsActive = false

	err = service.repository.Update(ctx, roomID, rooms[0])
	if err != nil {
		return err
	}

	return nil
}
