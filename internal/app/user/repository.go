package user

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
	repositoryName = "user.repository"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) error
	Get(ctx context.Context, userSearch entities.UserSearch) ([]entities.User, error)
	FindOne(ctx context.Context, userSearch entities.UserSearch) (entities.User, error)
	Update(ctx context.Context, userID string, user entities.User) error
}

type userRepository struct {
	config         config.Config
	mongodb        mongodb.MongoDBier
	logs           logger.Logger
	collectionName string
}

func NewUserRepository(cfg config.Config, mongoDBier mongodb.MongoDBier, logger logger.Logger) UserRepository {
	return &userRepository{
		config:         cfg,
		mongodb:        mongoDBier,
		logs:           logger,
		collectionName: cfg.MongoDB.Collections.Users,
	}
}

func (repository *userRepository) Create(ctx context.Context, user entities.User) error {
	userDTO, err := entities.CreateUserDTOFromUserEntity(user)
	if err != nil {
		repository.logs.Error(err.Error(), fmt.Sprintf("%s.%s", repositoryName, "Set"))
		return err
	}

	_, err = repository.mongodb.Collection(repository.collectionName).InsertOne(ctx, userDTO)
	if err != nil {
		repository.logs.Error(err.Error(), fmt.Sprintf("%s.%s", repositoryName, "Set"))
		return err
	}

	return nil
}

func (repository *userRepository) Get(ctx context.Context, search entities.UserSearch) ([]entities.User, error) {
	var users []entities.User
	collection := repository.mongodb.Collection(repository.collectionName)
	filter := createFilter(search)
	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		return users, err
	}

	defer func() {
		errClose := cursor.Close(ctx)
		if errClose != nil {
			repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
		}
	}()

	for cursor.Next(ctx) {
		userDTO := new(entities.UserDTO)
		err = cursor.Decode(userDTO)
		if err != nil {
			repository.logs.Error(str.ErrorConcat(err, repositoryName, "Get"))
			return users, err
		}

		users = append(users, entities.CreateUserEntityFromUserDTO(*userDTO))
	}

	return users, nil
}

func (repository *userRepository) FindOne(ctx context.Context, userSearch entities.UserSearch) (entities.User, error) {
	var user entities.User
	var userDTO entities.UserDTO
	collection := repository.mongodb.Collection(repository.collectionName)
	result := collection.FindOne(ctx, createFilter(userSearch))
	if result.Err() != nil {
		if result.Err().Error() == mongodb.NoResultsOnFind {
			err := exceptions.NewNotFoundException(fmt.Sprintf("user with username: %s not found", userSearch.Username))
			repository.logs.Warn(str.ErrorConcat(err, repositoryName, "FindOne"))
			return user, err
		}

		repository.logs.Error(str.ErrorConcat(result.Err(), repositoryName, "FindOne"))
		return user, result.Err()
	}

	err := result.Decode(&userDTO)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(result.Err(), repositoryName, "FindOne"))
		return user, err
	}

	user = entities.CreateUserEntityFromUserDTO(userDTO)

	return user, nil
}

func (repository *userRepository) Update(ctx context.Context, userID string, user entities.User) error {
	foundID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	Collection := repository.mongodb.Collection(repository.collectionName)
	filter := bson.M{entities.UserIDField: foundID}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				primitive.E{Key: entities.UsernameField, Value: user.Username},
				primitive.E{Key: entities.UserIsActiveNameField, Value: user.IsActive},
			},
		},
	}

	result, err := Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	if result.MatchedCount == 0 {
		err = exceptions.NewNotFoundException(fmt.Sprintf("user with UserID:%s not found", userID))
		repository.logs.Error(str.ErrorConcat(err, repositoryName, "Update"))
		return err
	}

	return nil
}

func createFilter(search entities.UserSearch) bson.D {
	var filter bson.D
	if !str.IsEmpty(search.ID) {
		id, _ := primitive.ObjectIDFromHex(search.ID)
		filter = append(filter, bson.E{Key: entities.UserIDField, Value: id})
	}

	if !str.IsEmpty(search.Username) {
		filter = append(filter, bson.E{Key: entities.UsernameField, Value: search.Username})
	}

	if search.IsActive != nil {
		filter = append(filter, bson.E{Key: entities.UserIsActiveNameField, Value: search.IsActive})
	}

	if len(filter) == 0 {
		return bson.D{
			{Key: entities.UserIDField, Value: bson.M{"$exists": true}},
			{Key: entities.UsernameField, Value: bson.M{"$exists": true}},
			{Key: entities.UserIsActiveNameField, Value: bson.M{"$exists": true}},
		}
	}

	return filter
}
