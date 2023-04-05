package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/equimper/twitter/config"
	"github.com/equimper/twitter/domain"
	"github.com/equimper/twitter/graph"
	"github.com/equimper/twitter/jwt"
	"github.com/equimper/twitter/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx := context.Background()

	config.LoadEnv(".env")

	conf := config.New()

	db := postgres.New(ctx, conf)

	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Timeout(time.Second * 60))

	// REPOS
	userRepo := postgres.NewUserRepo(db)
	tweetRepo := postgres.NewTweetRepo(db)

	// SERVICES
	authTokenService := jwt.NewTokenService(conf)
	authService := domain.NewAuthService(userRepo, authTokenService)
	tweetService := domain.NewTweetService(tweetRepo)
	userService := domain.NewUserService(userRepo)

	router.Use(graph.DataloaderMiddleware(
		&graph.Repos{
			UserRepo: userRepo,
		},
	))

	router.Use(authMiddleware(authTokenService))
	router.Handle("/", playground.Handler("Twitter clone", "/query"))
	router.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					AuthService:  authService,
					TweetService: tweetService,
					UserService:  userService,
				},
			},
		),
	))

	log.Fatal(http.ListenAndServe(":8080", router))
}
