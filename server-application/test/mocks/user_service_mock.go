package mocks

import (
	"context"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return new(UserServiceMock)
}

func (m *UserServiceMock) Create(ctx context.Context, user entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserServiceMock) Login(ctx context.Context, user entities.User) (entities.UserLoginResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entities.UserLoginResponse), args.Error(1)
}

func (m *UserServiceMock) Get(ctx context.Context, search entities.UserSearch) (entities.UsersSearchResponse, error) {
	args := m.Called(ctx, search)
	return args.Get(0).(entities.UsersSearchResponse), args.Error(1)
}

func (m *UserServiceMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
