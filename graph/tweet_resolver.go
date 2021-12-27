package graph

import (
	"context"
	"errors"

	"github.com/equimper/twitter"
)

func mapTweet(t twitter.Tweet) *Tweet {
	return &Tweet{
		ID:        t.ID,
		Body:      t.Body,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
	}
}

func mapTweets(tweets []twitter.Tweet) []*Tweet {
	tt := make([]*Tweet, len(tweets))

	for i, t := range tweets {
		tt[i] = mapTweet(t)
	}

	return tt
}

func (q *queryResolver) Tweets(ctx context.Context) ([]*Tweet, error) {
	tweets, err := q.TweetService.All(ctx)
	if err != nil {
		return nil, err
	}

	return mapTweets(tweets), nil
}

func (m *mutationResolver) CreateTweet(ctx context.Context, input CreateTweetInput) (*Tweet, error) {
	tweet, err := m.TweetService.Create(ctx, twitter.CreateTweetInput{
		Body: input.Body,
	})
	if err != nil {
		switch {
		case errors.Is(err, twitter.ErrUnauthenticated):
			return nil, buildUnauthenticatedError(ctx, err)
		case errors.Is(err, twitter.ErrValidation):
			return nil, buildBadRequestError(ctx, err)
		default:
			return nil, err
		}
	}

	return mapTweet(tweet), nil
}
