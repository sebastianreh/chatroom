package user_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebastianreh/chatroom/internal/app/user"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/internal/entities/exceptions"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_UserService_Create(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("create user successfully", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			ID:       "id123",
			Username: "User1",
			Password: "newPassword",
			IsActive: true,
		}
		userSearch := entities.UserSearch{Username: request.Username}

		repositoryMock.On("Get", ctx, userSearch).Return([]entities.User{}, nil)
		repositoryMock.On("Create", ctx, request).Return(nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Create(ctx, request)

		assert.NoError(t, err)
	})

	t.Run("user already exists", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			ID:       "id456",
			Username: "User2",
			Password: "password123",
			IsActive: true,
		}
		userSearch := entities.UserSearch{Username: request.Username}

		existingUsers := []entities.User{request}
		repositoryMock.On("Get", ctx, userSearch).Return(existingUsers, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Create(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewDuplicatedException(fmt.Sprintf("user '%s' already exist", request.Username)), err)
	})

	t.Run("error getting users from repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			ID:       "id789",
			Username: "User3",
			Password: "password456",
			IsActive: true,
		}
		userSearch := entities.UserSearch{Username: request.Username}

		expectedErr := errors.New("database error")
		repositoryMock.On("Get", ctx, userSearch).Return([]entities.User{}, expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Create(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("error creating user in repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			ID:       "id012",
			Username: "User4",
			Password: "password789",
			IsActive: true,
		}
		userSearch := entities.UserSearch{Username: request.Username}

		expectedErr := errors.New("database error")
		repositoryMock.On("Get", ctx, userSearch).Return([]entities.User{}, nil)
		repositoryMock.On("Create", ctx, request).Return(expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Create(ctx, request)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
	})
}

func Test_UserService_Login(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("successful login", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			Username: "User1",
			Password: "password123",
		}

		hashedPassword, _ := entities.HashPassword(request.Password)
		userFound := entities.User{
			ID:       "id123",
			Username: "User1",
			Password: hashedPassword,
			IsActive: true,
		}

		repositoryMock.On("FindOne", ctx, entities.UserSearch{Username: request.Username}).Return(userFound, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		resp, err := service.Login(ctx, request)

		assert.NoError(t, err)
		assert.Equal(t, userFound.ID, resp.ID)
	})

	t.Run("user does not exist", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			Username: "UserNotExist",
			Password: "password123",
		}

		emptyUser := entities.User{}
		repositoryMock.On("FindOne", ctx, entities.UserSearch{Username: request.Username}).Return(emptyUser, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewNotFoundException(fmt.Sprintf("user '%s' does not exist", request.Username)), err)
	})

	t.Run("user's account is not active", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			Username: "User2",
			Password: "password123",
		}

		userFound := entities.User{
			ID:       "id456",
			Username: "User2",
			Password: "hashedPassword123", // This can be any hashed password, since the user is inactive and the password won't be checked
			IsActive: false,
		}

		repositoryMock.On("FindOne", ctx, entities.UserSearch{Username: request.Username}).Return(userFound, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewUnauthorizedException(fmt.Sprintf("user '%s' account is not active", request.Username)), err)
	})

	t.Run("user password does not match", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			Username: "User3",
			Password: "wrongPassword",
		}

		hashedPassword, _ := entities.HashPassword("correctPassword")
		userFound := entities.User{
			ID:       "id789",
			Username: "User3",
			Password: hashedPassword,
			IsActive: true,
		}

		repositoryMock.On("FindOne", ctx, entities.UserSearch{Username: request.Username}).Return(userFound, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewUnauthorizedException(fmt.Sprintf("user %s credentials don't match", request.Username)), err)
	})

	t.Run("error fetching user from repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		request := entities.User{
			Username: "User4",
			Password: "password456",
		}

		expectedErr := errors.New("database error")
		repositoryMock.On("FindOne", ctx, entities.UserSearch{Username: request.Username}).Return(entities.User{}, expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Login(ctx, request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_UserService_Get(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("successful retrieval of users", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		search := entities.UserSearch{
			Username: "User",
		}

		users := []entities.User{
			{
				ID:       "id123",
				Username: "User1",
				IsActive: true,
			},
			{
				ID:       "id456",
				Username: "User2",
				IsActive: true,
			},
		}

		repositoryMock.On("Get", ctx, search).Return(users, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		resp, err := service.Get(ctx, search)

		assert.NoError(t, err)
		assert.Equal(t, len(users), len(resp.Users))
	})

	t.Run("no users found with given filter", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		search := entities.UserSearch{
			Username: "NonExistentUser",
		}

		repositoryMock.On("Get", ctx, search).Return([]entities.User{}, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Get(ctx, search)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewNotFoundException("users by filter not found"), err)
	})

	t.Run("error fetching users from repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		search := entities.UserSearch{
			Username: "UserWithError",
		}

		expectedErr := errors.New("database error")
		repositoryMock.On("Get", ctx, search).Return([]entities.User{}, expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		_, err := service.Get(ctx, search)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func Test_UserService_Delete(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("successful deletion of user", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		userID := "id123"

		userFound := entities.User{
			ID:       userID,
			Username: "User1",
			IsActive: true,
		}

		repositoryMock.On("FindOne", ctx, entities.UserSearch{ID: userID}).Return(userFound, nil)
		repositoryMock.On("Update", ctx, userID, mock.MatchedBy(func(user entities.User) bool {
			return !user.IsActive
		})).Return(nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Delete(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("user with given ID does not exist", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		userID := "idNotExist"

		repositoryMock.On("FindOne", ctx, entities.UserSearch{ID: userID}).Return(entities.User{}, nil)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Delete(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, exceptions.NewNotFoundException(fmt.Sprintf("user with UserID '%s' does not exist", userID)), err)
	})

	t.Run("error fetching user from repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		userID := "idWithError"

		expectedErr := errors.New("database error")
		repositoryMock.On("FindOne", ctx, entities.UserSearch{ID: userID}).Return(entities.User{}, expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Delete(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("error updating user in repository", func(t *testing.T) {
		repositoryMock := mocks.NewUserRepositoryMock()
		ctx := context.TODO()

		userID := "idUpdateError"

		userFound := entities.User{
			ID:       userID,
			Username: "UserError",
			IsActive: true,
		}

		expectedErr := errors.New("update error")
		repositoryMock.On("FindOne", ctx, entities.UserSearch{ID: userID}).Return(userFound, nil)
		repositoryMock.On("Update", ctx, userID, mock.Anything).Return(expectedErr)

		service := user.NewUserService(configs, repositoryMock, logs)

		err := service.Delete(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
