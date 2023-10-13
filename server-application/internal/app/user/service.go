package user

import (
	"context"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/internal/entities/exceptions"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
)

const serviceName = "user.service"

type UserService interface {
	Create(ctx context.Context, user entities.User) error
	Login(ctx context.Context, user entities.User) (entities.UserLoginResponse, error)
	Get(ctx context.Context, search entities.UserSearch) (entities.UsersSearchResponse, error)
	Delete(ctx context.Context, userID string) error
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

func (service *userService) Create(ctx context.Context, user entities.User) error {
	users, err := service.repository.Get(ctx, entities.UserSearch{Username: user.Username})
	if err != nil {
		return err
	}

	if len(users) > 0 {
		err = exceptions.NewDuplicatedException(fmt.Sprintf("user '%s' already exist", user.Username))
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return err
	}

	err = service.repository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) Login(ctx context.Context, user entities.User) (entities.UserLoginResponse, error) {
	var response entities.UserLoginResponse
	userFound, err := service.repository.FindOne(ctx, entities.UserSearch{Username: user.Username})
	if err != nil {
		return response, err
	}

	if userFound.IsEmpty() {
		err = exceptions.NewNotFoundException(fmt.Sprintf("user '%s' does not exist", user.Username))
		service.logs.Warn(str.ErrorConcat(err, serviceName, "Login"))
		return response, err
	}

	if userFound.IsActive == false {
		err = exceptions.NewUnauthorizedException(fmt.Sprintf("user '%s' account is not active", user.Username))
		service.logs.Warn(str.ErrorConcat(err, serviceName, "Login"))
		return response, err
	}

	err = entities.CompareHashAndPassword(userFound.Password, user.Password)
	if err != nil {
		err = exceptions.NewUnauthorizedException(fmt.Sprintf("user %s credentials don't match", userFound.Username))
		service.logs.Warn(str.ErrorConcat(err, serviceName, "Login"))
		return response, err
	}

	response.ID = userFound.ID

	return response, nil
}

func (service *userService) Get(ctx context.Context, search entities.UserSearch) (entities.UsersSearchResponse, error) {
	var usersSearchResponse entities.UsersSearchResponse
	usersFound, err := service.repository.Get(ctx, search)
	if err != nil {
		return usersSearchResponse, err
	}

	if len(usersFound) == 0 {
		err = exceptions.NewNotFoundException("users by filter not found")
		service.logs.Warn(str.ErrorConcat(err, repositoryName, "Get"))
		return usersSearchResponse, err
	}

	usersSearchResponse = entities.CreateUsersSearchResponseFromSearch(usersFound)

	return usersSearchResponse, nil
}

func (service *userService) Delete(ctx context.Context, userID string) error {
	userFound, err := service.repository.FindOne(ctx, entities.UserSearch{ID: userID})
	if err != nil {
		return err
	}

	if userFound.IsEmpty() {
		err = exceptions.NewNotFoundException(fmt.Sprintf("user with UserID '%s' does not exist", userID))
		service.logs.Error(str.ErrorConcat(err, serviceName, "Set"))
		return err
	}

	userFound.IsActive = false
	err = service.repository.Update(ctx, userID, userFound)
	if err != nil {
		return err
	}

	return nil
}
