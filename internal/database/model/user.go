// Package model contains all databank models
package model

import "github.com/topvennie/spotify_organizer/pkg/sqlc"

type User struct {
	ID    int
	Name  string
	Email string
	UID   string
}

func UserModel(user sqlc.User) *User {
	return &User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		UID:   user.Uid,
	}
}
