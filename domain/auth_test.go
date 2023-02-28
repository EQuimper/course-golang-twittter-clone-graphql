package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{
		Username:        "bob",
		Email:           "bob@example.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(twitter.User{
				ID:       "123",
				Username: validInput.Username,
				Email:    validInput.Email,
				Password: validInput.Password,
			}, nil)

		service := NewAuthService(userRepo)

		res, err := service.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.AccessToken)
		require.NotEmpty(t, res.User.ID)
		require.NotEmpty(t, res.User.Email)
		require.NotEmpty(t, res.User.Username)
		require.NotEmpty(t, res.User.Password)

		userRepo.AssertExpectations(t)
	})

	t.Run("username taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, nil)

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrUsernameTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})

	t.Run("email taken", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(twitter.User{}, nil)

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.ErrorIs(t, err, twitter.ErrEmailTaken)

		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})

	t.Run("create error", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		userRepo.On("GetByUsername", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("GetByEmail", mock.Anything, mock.Anything).
			Return(twitter.User{}, twitter.ErrNotFound)

		userRepo.On("Create", mock.Anything, mock.Anything).
			Return(twitter.User{}, errors.New("something"))

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, validInput)
		require.Error(t, err)

		userRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		ctx := context.Background()

		userRepo := &mocks.UserRepo{}

		service := NewAuthService(userRepo)

		_, err := service.Register(ctx, twitter.RegisterInput{})
		require.ErrorIs(t, err, twitter.ErrValidation)

		userRepo.AssertNotCalled(t, "GetByUsername")
		userRepo.AssertNotCalled(t, "GetByEmail")
		userRepo.AssertNotCalled(t, "Create")

		userRepo.AssertExpectations(t)
	})
}
