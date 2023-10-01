package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	RoomIDField       = "_id"
	RoomNameField     = "name"
	IsActiveNameField = "is_active"
)

type Room struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type RoomsGetResponse struct {
	Rooms []Room `json:"rooms"`
}

type RoomDTO struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name"  bson:"name"`
	IsActive bool               `json:"is_active"  bson:"is_active"`
}

type RoomSearch struct {
	ID       string `query:"id" bson:"_id"`
	Name     string `query:"name"  bson:"name"`
	IsActive *bool  `query:"is_active"  bson:"is_active"`
}

func CreateRoomDTOFromRequest(request Room) RoomDTO {
	return RoomDTO{
		ID:       primitive.NewObjectID(),
		Name:     request.Name,
		IsActive: true,
	}
}

func CreateRoomFromRoomDTO(DTO RoomDTO) Room {
	return Room{
		ID:       DTO.ID.Hex(),
		Name:     DTO.Name,
		IsActive: DTO.IsActive,
	}
}
