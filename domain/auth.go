package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/equimper/twitter"
	"golang.org/x/crypto/bcrypt"
)

var passwordCost = bcrypt.DefaultCost

type AuthService struct {
	AuthTokenService twitter.AuthTokenService
	UserRepo         twitter.UserRepo
}

func NewAuthService(ur twitter.UserRepo, ats twitter.AuthTokenService) *AuthService {
	return &AuthService{
		AuthTokenService: ats,
		UserRepo:         ur,
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
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), passwordCost)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(hashPassword)

	// create the user
	user, err = as.UserRepo.Create(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, fmt.Errorf("error creating user: %v", err)
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, twitter.ErrGenAccessToken
	}

	// return accessToken and user
	return twitter.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}

func (as *AuthService) Login(ctx context.Context, input twitter.LoginInput) (twitter.AuthResponse, error) {
	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.AuthResponse{}, err
	}

	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, twitter.ErrNotFound):
			return twitter.AuthResponse{}, twitter.ErrBadCredentials
		default:
			return twitter.AuthResponse{}, err
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return twitter.AuthResponse{}, twitter.ErrBadCredentials
	}

	accessToken, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return twitter.AuthResponse{}, twitter.ErrGenAccessToken
	}

	return twitter.AuthResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}
