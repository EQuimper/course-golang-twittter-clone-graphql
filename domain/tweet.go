package domain

import (
	"context"

	"github.com/equimper/twitter"
)

type TweetService struct {
	TweetRepo twitter.TweetRepo
}

func NewTweetService(tr twitter.TweetRepo) *TweetService {
	return &TweetService{
		TweetRepo: tr,
	}
}

func (ts *TweetService) All(ctx context.Context) ([]twitter.Tweet, error) {
	panic("not implemented") // TODO: Implement
}

func (ts *TweetService) Create(ctx context.Context, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	_, err := twitter.GetUserIDFromContext(ctx)
	if err != nil {
		return twitter.Tweet{}, twitter.ErrUnauthenticated
	}

	return twitter.Tweet{}, nil
}

func (ts *TweetService) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	panic("not implemented") // TODO: Implement
}
