package user_test

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sebastianreh/chatroom/cmd/httpserver"
	"github.com/sebastianreh/chatroom/cmd/httpserver/resterror"
	"github.com/sebastianreh/chatroom/internal/app/user"
	"github.com/sebastianreh/chatroom/internal/config"
	"github.com/sebastianreh/chatroom/internal/container"
	"github.com/sebastianreh/chatroom/internal/entities"
	"github.com/sebastianreh/chatroom/pkg/logger"
	"github.com/sebastianreh/chatroom/test/mocks"
	"github.com/stretchr/testify/assert"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(method, target string, body *strings.Reader) (echo.Context, *httptest.ResponseRecorder) {
	mockServer := httpserver.NewServer(container.Dependencies{})

	request := httptest.NewRequest(method, target, body)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	mockServer.Server.HTTPErrorHandler = httpserver.HTTPErrorHandler
	ctx := mockServer.NewServerContext(request, w)

	return ctx, w
}

func setPathAndParams(ctx echo.Context, path, names, values string) {
	ctx.SetPath(path)
	ctx.SetParamNames(names)
	ctx.SetParamValues(values)
}

func Test_UserHandler_Create(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("create user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		request := entities.User{
			ID:       "id123",
			Username: "User1",
			Password: "newPassword",
			IsActive: true,
		}

		body, _ := json.Marshal(request)
		context, _ := setup(http.MethodGet, "/", strings.NewReader(string(body)))
		serviceMock.On("Create", context.Request().Context(), request).Return(nil)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Create(context)

		assert.NoError(t, err)
	})

	t.Run("service create returns error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		request := entities.User{
			ID:       "id123",
			Username: "User1",
			Password: "newPassword",
			IsActive: true,
		}

		expectedError := resterror.NewBadRequestError("repository error")
		body, _ := json.Marshal(request)
		context, recorder := setup(http.MethodGet, "/", strings.NewReader(string(body)))
		serviceMock.On("Create", context.Request().Context(), request).Return(expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Create(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, context.Response().Status)
		restError, _ := resterror.NewRestErrorFromBytes(recorder.Body.Bytes())
		assert.Equal(t, expectedError, restError)
	})
}

func Test_UserHandler_Login(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("login user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		request := entities.User{
			Username: "User1",
			Password: "password123",
		}

		response := entities.UserLoginResponse{
			ID: "id123",
		}

		body, _ := json.Marshal(request)
		context, _ := setup(http.MethodPost, "/login", strings.NewReader(string(body)))
		serviceMock.On("Login", context.Request().Context(), request).Return(response, nil)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Login(context)

		assert.NoError(t, err)
	})

	t.Run("bind error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		body := strings.NewReader("invalid json")
		context, recorder := setup(http.MethodPost, "/login", body)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Login(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("service login error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		request := entities.User{
			Username: "User1",
			Password: "password123",
		}

		expectedError := resterror.NewBadRequestError("Service error")
		body, _ := json.Marshal(request)
		context, recorder := setup(http.MethodPost, "/login", strings.NewReader(string(body)))
		serviceMock.On("Login", context.Request().Context(), request).Return(entities.UserLoginResponse{}, expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Login(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("invalid login credentials", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		request := entities.User{
			Username: "User1",
			Password: "wrongpassword",
		}

		expectedError := resterror.NewUnauthorizedError("Invalid credentials")
		body, _ := json.Marshal(request)
		context, recorder := setup(http.MethodPost, "/login", strings.NewReader(string(body)))
		serviceMock.On("Login", context.Request().Context(), request).Return(entities.UserLoginResponse{}, expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Login(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

func Test_UserHandler_Get(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("get user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userSearch := entities.UserSearch{
			Username: "User1",
		}

		response := entities.UsersSearchResponse{
			Users: []entities.UserSearchResponse{
				{
					ID:       "id123",
					Username: "User1",
					IsActive: true,
				},
			},
		}

		body, _ := json.Marshal(userSearch)
		context, _ := setup(http.MethodGet, "/get", strings.NewReader(string(body)))
		serviceMock.On("Get", context.Request().Context(), userSearch).Return(response, nil)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Get(context)

		assert.NoError(t, err)
	})

	t.Run("bind error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		body := strings.NewReader("invalid json")
		context, recorder := setup(http.MethodGet, "/get", body)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Get(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("service get error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userSearch := entities.UserSearch{
			Username: "User1",
		}

		expectedError := resterror.NewInternalServerError("Service error", errors.New("error"))
		body, _ := json.Marshal(userSearch)
		context, recorder := setup(http.MethodGet, "/get", strings.NewReader(string(body)))
		serviceMock.On("Get", context.Request().Context(), userSearch).Return(entities.UsersSearchResponse{}, expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Get(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("empty search result", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userSearch := entities.UserSearch{
			Username: "NonexistentUser",
		}

		body, _ := json.Marshal(userSearch)
		context, recorder := setup(http.MethodGet, "/get", strings.NewReader(string(body)))
		serviceMock.On("Get", context.Request().Context(), userSearch).Return(entities.UsersSearchResponse{}, nil)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Get(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func Test_UserHandler_Delete(t *testing.T) {
	logs := logger.NewLogger()
	configs := config.NewConfig()

	t.Run("delete user successfully", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userID := "id123"

		context, _ := setup(http.MethodDelete, "/delete/"+userID, strings.NewReader(""))
		setPathAndParams(context, "/delete/:id", "id", userID)
		serviceMock.On("Delete", context.Request().Context(), userID).Return(nil)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Delete(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, context.Response().Status)
	})

	t.Run("empty user id", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()

		context, recorder := setup(http.MethodDelete, "/delete/", strings.NewReader(""))
		setPathAndParams(context, "/delete/:id", "id", "")
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Delete(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("service delete error", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userID := "id123"

		expectedError := resterror.NewInternalServerError("Service error", errors.New(""))
		context, recorder := setup(http.MethodDelete, "/delete/"+userID, &strings.Reader{})
		setPathAndParams(context, "/delete/:id", "id", userID)
		serviceMock.On("Delete", context.Request().Context(), userID).Return(expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Delete(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		serviceMock := mocks.NewUserServiceMock()
		userID := "nonexistentID"

		expectedError := resterror.NewNotFoundError("User not found")
		context, recorder := setup(http.MethodDelete, "/delete/"+userID, &strings.Reader{})
		setPathAndParams(context, "/delete/:id", "id", userID)
		serviceMock.On("Delete", context.Request().Context(), userID).Return(expectedError)
		handler := user.NewUserHandler(configs, serviceMock, logs)

		err := handler.Delete(context)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
