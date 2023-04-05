package postgres

import (
	"context"
	"fmt"

	"github.com/equimper/twitter"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type TweetRepo struct {
	DB *DB
}

func NewTweetRepo(db *DB) *TweetRepo {
	return &TweetRepo{
		DB: db,
	}
}

func (tr *TweetRepo) All(ctx context.Context) ([]twitter.Tweet, error) {
	return getAllTweets(ctx, tr.DB.Pool)
}

func getAllTweets(ctx context.Context, q pgxscan.Querier) ([]twitter.Tweet, error) {
	query := `SELECT * FROM tweets ORDER BY created_at DESC;`

	var tweets []twitter.Tweet

	if err := pgxscan.Select(ctx, q, &tweets, query); err != nil {
		return nil, fmt.Errorf("error get all tweets %+v", err)
	}

	return tweets, nil
}

func (tr *TweetRepo) Create(ctx context.Context, tweet twitter.Tweet) (twitter.Tweet, error) {
	tx, err := tr.DB.Pool.Begin(ctx)
	if err != nil {
		return twitter.Tweet{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	tweet, err = createTweet(ctx, tx, tweet)
	if err != nil {
		return twitter.Tweet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error commiting: %v", err)
	}

	return tweet, nil
}

func createTweet(ctx context.Context, tx pgx.Tx, tweet twitter.Tweet) (twitter.Tweet, error) {
	query := `INSERT INTO tweets (body, user_id, parent_id) VALUES ($1, $2, $3) RETURNING *;`

	t := twitter.Tweet{}

	if err := pgxscan.Get(ctx, tx, &t, query, tweet.Body, tweet.UserID, tweet.ParentID); err != nil {
		return twitter.Tweet{}, fmt.Errorf("error insert: %v", err)
	}

	return t, nil
}

func (tr *TweetRepo) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	return getTweetByID(ctx, tr.DB.Pool, id)
}

func getTweetByID(ctx context.Context, q pgxscan.Querier, id string) (twitter.Tweet, error) {
	query := `SELECT * FROM tweets WHERE id = $1 LIMIT 1;`

	t := twitter.Tweet{}

	if err := pgxscan.Get(ctx, q, &t, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.Tweet{}, twitter.ErrNotFound
		}

		return twitter.Tweet{}, fmt.Errorf("error get tweet: %+v", err)
	}

	return t, nil
}

func (tr *TweetRepo) Delete(ctx context.Context, id string) error {
	tx, err := tr.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := deleteTweet(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting: %v", err)
	}

	return nil
}

func deleteTweet(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM tweets WHERE id = $1;`

	if _, err := tx.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("error insert: %v", err)
	}

	return nil
}
