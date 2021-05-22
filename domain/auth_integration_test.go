// +build integration

package domain

import (
	"context"
	"testing"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/faker"
	"github.com/equimper/twitter/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAuthService_Register(t *testing.T) {
	validInput := twitter.RegisterInput{
		Username:        faker.Username(),
		Email:           faker.Email(),
		Password:        "password",
		ConfirmPassword: "password",
	}

	t.Run("can register a user", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		res, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		require.NotEmpty(t, res.User.ID)
		require.Equal(t, validInput.Email, res.User.Email)
		require.Equal(t, validInput.Username, res.User.Username)
		require.NotEqual(t, validInput.Password, res.User.Password)
	})

	t.Run("existing username", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter.RegisterInput{
			Username:        validInput.Username,
			Email:           faker.Email(),
			Password:        "password",
			ConfirmPassword: "password",
		})
		require.ErrorIs(t, err, twitter.ErrUsernameTaken)
	})

	t.Run("existing email", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := authService.Register(ctx, validInput)
		require.NoError(t, err)

		_, err = authService.Register(ctx, twitter.RegisterInput{
			Username:        faker.Username(),
			Email:           validInput.Email,
			Password:        "password",
			ConfirmPassword: "password",
		})
		require.ErrorIs(t, err, twitter.ErrEmailTaken)
	})
}
