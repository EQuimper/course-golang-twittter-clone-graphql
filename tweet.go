package twitter

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var (
	TweetMinLength = 2
	TweetMaxLength = 250
)

type CreateTweetInput struct {
	Body string
}

func (in *CreateTweetInput) Sanitize() {
	in.Body = strings.TrimSpace(in.Body)
}

func (in CreateTweetInput) Validate() error {
	if len(in.Body) < TweetMinLength {
		return fmt.Errorf("%w: body not long enough, (%d) characters at least", ErrValidation, TweetMinLength)
	}

	if len(in.Body) > TweetMaxLength {
		return fmt.Errorf("%w: body too long, (%d) characters at max", ErrValidation, TweetMaxLength)
	}

	return nil
}

type Tweet struct {
	ID        string
	Body      string
	UserID    string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Tweet) CanDelete(user User) bool {
	return t.UserID == user.ID
}

type TweetService interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, input CreateTweetInput) (Tweet, error)
	CreateReply(ctx context.Context, parentID string, input CreateTweetInput) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}

type TweetRepo interface {
	All(ctx context.Context) ([]Tweet, error)
	Create(ctx context.Context, tweet Tweet) (Tweet, error)
	GetByID(ctx context.Context, id string) (Tweet, error)
	Delete(ctx context.Context, id string) error
}
