package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username"  bson:"username"`
	Password string             `json:"password" bson:"password"`
	IsActive bool               `json:"is_active"  bson:"is_active"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func CreateUserDTOFromRequest(request User) (*UserDTO, error) {
	hashedPassword, err := HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	userDTO := &UserDTO{
		ID:       primitive.NewObjectID(),
		Username: request.Username,
		Password: hashedPassword,
		IsActive: true,
	}

	return userDTO, nil
}
