//go:build integration
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

func TestIntegrationTweetService_All(t *testing.T) {
	t.Run("return all tweets", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)

		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)
		test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweets, err := tweetService.All(ctx)
		require.NoError(t, err)

		require.Len(t, tweets, 3)
	})
}

func TestIntegrationTweetService_GetByID(t *testing.T) {
	t.Run("can get a tweet by id", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		user := test_helpers.CreateUser(ctx, t, userRepo)
		existingTweet := test_helpers.CreateTweet(ctx, t, tweetRepo, user.ID)

		tweet, err := tweetService.GetByID(ctx, existingTweet.ID)
		require.NoError(t, err)

		require.Equal(t, existingTweet.ID, tweet.ID, "tweet.ID")
		require.Equal(t, existingTweet.Body, tweet.Body, "tweet.Body")
	})

	t.Run("return error not found if the tweet doesn't exist", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, faker.UUID())
		require.ErrorIs(t, err, twitter.ErrNotFound)
	})

	t.Run("return error invalid uuid", func(t *testing.T) {
		ctx := context.Background()

		defer test_helpers.TeardownDB(ctx, t, db)

		_, err := tweetService.GetByID(ctx, "123")
		require.ErrorIs(t, err, twitter.ErrInvalidUUID)
	})
}
