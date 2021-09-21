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

func TestIntegrationTweetService_Create(t *testing.T) {
	t.Run("not auth user cannot create a tweet", func(t *testing.T) {
		ctx := context.Background()

		_, err := tweetService.Create(ctx, twitter.CreateTweetInput{
			Body: "hello",
		})

		require.ErrorIs(t, err, twitter.ErrUnauthenticated)
	})

	t.Run("can create a tweet", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		currentUser := test_helpers.CreateUser(ctx, t, userRepo)

		ctx = test_helpers.LoginUser(ctx, t, currentUser)

		input := twitter.CreateTweetInput{
			Body: faker.RandStr(100),
		}

		tweet, err := tweetService.Create(ctx, input)
		require.NoError(t, err)

		require.NotEmpty(t, tweet.ID, "tweet.ID")
		require.Equal(t, input.Body, tweet.Body, "tweet.Body")
		require.Equal(t, currentUser.ID, tweet.UserID, "tweet.UserID")
		require.NotEmpty(t, tweet.CreatedAt, "tweet.CreatedAt")
	})
}
