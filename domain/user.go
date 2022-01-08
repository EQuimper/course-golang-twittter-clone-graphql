package domain

import (
	"context"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/uuid"
)

type UserService struct {
	UserRepo twitter.UserRepo
}

func NewUserService(ur twitter.UserRepo) *UserService {
	return &UserService{
		UserRepo: ur,
	}
}

func (u *UserService) GetByID(ctx context.Context, id string) (twitter.User, error) {
	if !uuid.Validate(id) {
		return twitter.User{}, twitter.ErrInvalidUUID
	}

	return u.UserRepo.GetByID(ctx, id)
}
