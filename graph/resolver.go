package graph

import (
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/equimper/twitter"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	AuthService  twitter.AuthService
	TweetService twitter.TweetService
	UserService  twitter.UserService
}

type queryResolver struct {
	*Resolver
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type tweetResolver struct {
	*Resolver
}

func (r *Resolver) Tweet() TweetResolver {
	return &tweetResolver{r}
}

func buildBadRequestError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusBadRequest,
		},
	}
}

func buildUnauthenticatedError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusUnauthorized,
		},
	}
}

func buildForbiddenError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildNotFoundError(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: err.Error(),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": http.StatusForbidden,
		},
	}
}

func buildError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, twitter.ErrForbidden):
		return buildForbiddenError(ctx, err)
	case errors.Is(err, twitter.ErrUnauthenticated):
		return buildUnauthenticatedError(ctx, err)
	case errors.Is(err, twitter.ErrValidation):
		return buildBadRequestError(ctx, err)
	case errors.Is(err, twitter.ErrNotFound):
		return buildNotFoundError(ctx, err)
	default:
		return err
	}
}
