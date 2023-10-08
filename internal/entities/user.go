package entities

import (
	str "github.com/sebastianreh/chatroom/pkg/strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserIDField           = "_id"
	UsernameField         = "username"
	UserIsActiveNameField = "is_active"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
	IsActive bool   `json:"is_active"`
}

type (
	UsersSearchResponse struct {
		Users []UserSearchResponse `json:"users"`
	}
	UserSearchResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		IsActive bool   `json:"is_active"`
	}
)

type UserDTO struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username"  bson:"username"`
	Password string             `json:"password" bson:"password"`
	IsActive bool               `json:"is_active"  bson:"is_active"`
}

type UserSearch struct {
	ID       string `query:"id" bson:"_id"`
	Username string `query:"username"  bson:"username"`
	IsActive *bool  `query:"is_active"  bson:"is_active"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func CompareHashAndPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func CreateUserDTOFromUserEntity(request User) (UserDTO, error) {
	var user UserDTO
	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		return user, err
	}

	userDTO := UserDTO{
		ID:       primitive.NewObjectID(),
		Username: request.Username,
		Password: hashedPassword,
		IsActive: true,
	}

	return userDTO, nil
}

func CreateUserEntityFromUserDTO(DTO UserDTO) User {
	return User{
		ID:       DTO.ID.Hex(),
		Username: DTO.Username,
		IsActive: DTO.IsActive,
	}
}

func (u User) IsEmpty() bool {
	if str.IsEmpty(u.Username) && str.IsEmpty(u.Password) {
		return true
	}
	return false
}

func CreateUsersSearchResponseFromSearch(usersFound []User) UsersSearchResponse {
	var usersSearchResponse UsersSearchResponse
	for _, user := range usersFound {
		usersSearchResponse.Users = append(usersSearchResponse.Users, UserSearchResponse{
			ID:       user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	return usersSearchResponse
}
