package domain

import (
	"context"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/uuid"
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
	return ts.TweetRepo.All(ctx)
}

func (ts *TweetService) Create(ctx context.Context, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	currentUserID, err := twitter.GetUserIDFromContext(ctx)
	if err != nil {
		return twitter.Tweet{}, twitter.ErrUnauthenticated
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.Tweet{}, err
	}

	tweet, err := ts.TweetRepo.Create(ctx, twitter.Tweet{
		Body:   input.Body,
		UserID: currentUserID,
	})
	if err != nil {
		return twitter.Tweet{}, err
	}

	return tweet, nil
}

func (ts *TweetService) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	if !uuid.Validate(id) {
		return twitter.Tweet{}, twitter.ErrInvalidUUID
	}

	return ts.TweetRepo.GetByID(ctx, id)
}

func (ts *TweetService) Delete(ctx context.Context, id string) error {
	currentUserID, err := twitter.GetUserIDFromContext(ctx)
	if err != nil {
		return twitter.ErrUnauthenticated
	}

	if !uuid.Validate(id) {
		return twitter.ErrInvalidUUID
	}

	tweet, err := ts.TweetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !tweet.CanDelete(twitter.User{ID: currentUserID}) {
		return twitter.ErrForbidden
	}

	return ts.TweetRepo.Delete(ctx, id)
}

func (ts *TweetService) CreateReply(ctx context.Context, parentID string, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	currentUserID, err := twitter.GetUserIDFromContext(ctx)
	if err != nil {
		return twitter.Tweet{}, twitter.ErrUnauthenticated
	}

	input.Sanitize()

	if err := input.Validate(); err != nil {
		return twitter.Tweet{}, err
	}

	if !uuid.Validate(parentID) {
		return twitter.Tweet{}, twitter.ErrInvalidUUID
	}

	if _, err := ts.TweetRepo.GetByID(ctx, parentID); err != nil {
		return twitter.Tweet{}, twitter.ErrNotFound
	}

	tweet, err := ts.TweetRepo.Create(ctx, twitter.Tweet{
		Body:     input.Body,
		UserID:   currentUserID,
		ParentID: &parentID,
	})
	if err != nil {
		return twitter.Tweet{}, err
	}

	return tweet, nil
}
