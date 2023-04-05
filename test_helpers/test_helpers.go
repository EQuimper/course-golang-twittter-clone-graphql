package test_helpers

import (
	"context"
	"testing"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/faker"
	"github.com/equimper/twitter/postgres"
	"github.com/stretchr/testify/require"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)
	require.NoError(t, err)
}

func CreateUser(ctx context.Context, t *testing.T, userRepo twitter.UserRepo) twitter.User {
	t.Helper()

	user, err := userRepo.Create(ctx, twitter.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password,
	})
	require.NoError(t, err)

	return user
}

func CreateTweet(ctx context.Context, t *testing.T, tweetRepo twitter.TweetRepo, forUser string) twitter.Tweet {
	t.Helper()

	tweet, err := tweetRepo.Create(ctx, twitter.Tweet{
		Body:   faker.RandStr(20),
		UserID: forUser,
	})
	require.NoError(t, err)

	return tweet
}

func LoginUser(ctx context.Context, t *testing.T, user twitter.User) context.Context {
	t.Helper()

	return twitter.PutUserIDIntoContext(ctx, user.ID)
}
