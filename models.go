package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/viniciuspra/rssagg/internal/db"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func dbUserToUser(dbUser db.User) User {
	return User{
		ID: dbUser.ID,
		Name: dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
