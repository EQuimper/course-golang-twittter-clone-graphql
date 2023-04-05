//go:generate go run github.com/vektah/dataloaden UserLoader string *twitter/graph.User

package graph

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/equimper/twitter"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserByID UserLoader
}

type Repos struct {
	UserRepo twitter.UserRepo
}

func DataloaderMiddleware(repos *Repos) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
				UserByID: UserLoader{
					wait:     1 * time.Millisecond,
					maxBatch: 100,
					fetch: func(ids []string) ([]*User, []error) {
						// []twitter.User
						users, err := repos.UserRepo.GetByIds(r.Context(), ids)
						if err != nil {
							return nil, []error{err}
						}

						userByID := map[string]*User{}

						for _, u := range users {
							userByID[u.ID] = mapUser(u)
						}

						result := make([]*User, len(ids))

						for i, id := range ids {
							user, ok := userByID[id]
							if !ok {
								return nil, []error{fmt.Errorf("user with id: %s is missing", id)}
							}

							result[i] = user
						}

						return result, nil
					},
				},
			})

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func DataloaderFor(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
