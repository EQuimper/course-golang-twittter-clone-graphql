package twitter

import "time"

type UserRepo interface {
}

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
