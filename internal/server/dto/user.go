package dto

import "github.com/topvennie/spotify_organizer/internal/database/model"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	UID   string `json:"-"`
}

func UserDTO(user model.User) User {
	return User(user)
}

func (u *User) ToModel() model.User {
	return model.User(*u)
}
