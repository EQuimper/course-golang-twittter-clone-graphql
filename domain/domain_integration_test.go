// +build integration

package domain

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/config"
	"github.com/equimper/twitter/jwt"
	"github.com/equimper/twitter/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	conf             *config.Config
	db               *postgres.DB
	authTokenService twitter.AuthTokenService
	authService      twitter.AuthService
	tweetService     twitter.TweetService
	userRepo         twitter.UserRepo
	tweetRepo        twitter.TweetRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	config.LoadEnv(".env.test")

	passwordCost = bcrypt.MinCost

	conf = config.New()

	db = postgres.New(ctx, conf)
	defer db.Close()

	if err := db.Drop(); err != nil {
		log.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo = postgres.NewUserRepo(db)
	tweetRepo = postgres.NewTweetRepo(db)

	authTokenService = jwt.NewTokenService(conf)

	authService = NewAuthService(userRepo, authTokenService)
	tweetService = NewTweetService(tweetRepo)

	os.Exit(m.Run())
}
