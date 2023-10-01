package room

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/internal/entities/exceptions"
	"github.com/sebastianreh/chatroom/pkg/logger"
)

const (
	serviceName = "room.service"
)

type RoomService interface {
	Create(ctx context.Context, room entities.Room) error
	Join(ctx context.Context) error
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
func (service *roomService) Create(ctx context.Context, room entities.Room) error {
	rooms, err := service.repository.Get(ctx, entities.RoomSearch{Name: room.Name})
	if err != nil {
		return err
	}

	if len(rooms) > 0 {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("room '%s' already exist", room.Name))
		service.logs.Error(err.Error(), fmt.Sprintf("%s.%s", serviceName, "Create"))
		return err
	}

	err = service.repository.Create(ctx, entities.CreateRoomDTOFromRequest(room))
	if err != nil {
		return err
	}

	return nil
}

func (service *roomService) Join(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (service *roomService) Get(ctx context.Context, search entities.RoomSearch) (entities.RoomsGetResponse, error) {
	var roomsResponse entities.RoomsGetResponse
	rooms, err := service.repository.Get(ctx, search)
	if err != nil {
		return roomsResponse, err
	}

	return entities.RoomsGetResponse{Rooms: rooms}, nil
}

func (service *roomService) Delete(ctx context.Context, roomID string) error {
	rooms, err := service.repository.Get(ctx, entities.RoomSearch{ID: roomID})
	if err != nil {
		return err
	}

	if len(rooms) != 1 {
		err = errors.New("found more than one room to delete")
		service.logs.Error(err.Error(), fmt.Sprintf("%s.%s", serviceName, "Get"))
		return err
	}

	rooms[0].IsActive = false

	err = service.repository.Update(ctx, roomID, rooms[0])
	if err != nil {
		return err
	}

	return nil
}
