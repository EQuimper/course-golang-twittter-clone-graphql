package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/equimper/twitter"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo twitter.UserRepo
}

func NewAuthService(ur twitter.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input twitter.RegisterInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	// check if username is already taken
	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrUsernameTaken
	}

	// check if email is already taken
	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, twitter.ErrNotFound) {
		return twitter.AuthResponse{}, twitter.ErrEmailTaken
	}

	user := twitter.User{
		Email:    input.Email,
		Username: input.Username,
	}

	// hash the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(hashPassword)

	// create the user
	user, err = as.UserRepo.Create(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	// return accessToken and user
	return twitter.AuthResponse{
		AccessToken: "a token",
		User:        user,
	}, nil
}
