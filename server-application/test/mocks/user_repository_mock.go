package mocks

import (
	"context"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return new(UserRepositoryMock)
}

func (m *UserRepositoryMock) Create(ctx context.Context, user entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindOne(ctx context.Context, userSearch entities.UserSearch) (entities.User, error) {
	args := m.Called(ctx, userSearch)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *UserRepositoryMock) Get(ctx context.Context, userSearch entities.UserSearch) ([]entities.User, error) {
	args := m.Called(ctx, userSearch)
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, userID string, user entities.User) error {
	args := m.Called(ctx, userID, user)
	return args.Error(0)
}
