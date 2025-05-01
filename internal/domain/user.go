package domain

import "time"

type User struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type UserCommand struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
