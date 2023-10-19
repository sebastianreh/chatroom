package room

import (
	"context"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/internal/entities/exceptions"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/pkg/mongodb"
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	repositoryName = "room.repository"
)

type RoomRepository interface {
	Create(ctx context.Context, room entities.Room) (string, error)
	Get(ctx context.Context, search entities.RoomSearch) ([]entities.Room, error)
	Update(ctx context.Context, roomID string, room entities.Room) error
}

type roomRepository struct {
	config         config.Config
	mongodb        mongodb.MongoDBier
	logs           logger.Logger
	collectionName string
}

func NewRoomRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logger.Logger) RoomRepository {
	return &roomRepository{
		config:         cfg,
		mongodb:        mongoDBier,
		logs:           logger,
		collectionName: cfg.MongoDB.Collections.Rooms,
	}
}

func (repository *roomRepository) Create(ctx context.Context, room entities.Room) (string, error) {
	roomDTO := entities.CreateRoomDTOFromEntity(room)
	_, err := repository.mongodb.Collection(repository.collectionName).InsertOne(ctx, roomDTO)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Set"))
		return str.Empty, err
	}

	return roomDTO.ID.Hex(), nil
}

func (repository *roomRepository) Get(ctx context.Context, search entities.RoomSearch) ([]entities.Room, error) {
	var rooms []entities.Room
	collection := repository.mongodb.Collection(repository.collectionName)
	filter := createFilter(search)
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		return rooms, err
	}

	defer func() {
		errClose := cursor.Close(ctx)
		if errClose != nil {
			repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		}
	}()

	for cursor.Next(ctx) {
		roomDTO := new(entities.RoomDTO)
		err = cursor.Decode(roomDTO)
		if err != nil {
			repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
			return rooms, err
		}

		rooms = append(rooms, entities.CreateRoomEntityFromRoomDTO(*roomDTO))
	}

	return rooms, nil
}

func (repository *roomRepository) Update(ctx context.Context, roomID string, room entities.Room) error {
	foundID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	Collection := repository.mongodb.Collection(repository.collectionName)
	filter := bson.M{entities.RoomIDField: foundID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: entities.RoomNameField, Value: room.Name},
				primitive.E{Key: entities.RoomIsActiveNameField, Value: room.IsActive},
			},
		},
	}

	result, err := Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = exceptions.NewNotFoundException(fmt.Sprintf("room with UserID:%s not found", roomID))
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	return nil
}

func createFilter(search entities.RoomSearch) bson.D {
	var filter bson.D
	if !str.IsEmpty(search.ID) {
		id, _ := primitive.ObjectIDFromHex(search.ID)
		filter = append(filter, bson.E{Key: entities.RoomIDField, Value: id})
	}

	if !str.IsEmpty(search.Name) {
		filter = append(filter, bson.E{Key: entities.RoomNameField, Value: search.Name})
	}

	if search.IsActive != nil {
		filter = append(filter, bson.E{Key: entities.RoomIsActiveNameField, Value: search.IsActive})
	}

	if len(filter) == 0 {
		return bson.D{}
	}

	return filter
}
