package twitter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("get user id from context", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, contextAuthIDKey, "123")

		userID, err := GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})

	t.Run("return error if no id", func(t *testing.T) {
		ctx := context.Background()

		_, err := GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, ErrNoUserIDInContext)
	})

	t.Run("return error if value is not a string", func(t *testing.T) {
		ctx := context.Background()

		ctx = context.WithValue(ctx, contextAuthIDKey, 123)

		_, err := GetUserIDFromContext(ctx)
		require.ErrorIs(t, err, ErrNoUserIDInContext)

	})
}

func TestPutUserIDIntoContext(t *testing.T) {
	t.Run("add user id into context", func(t *testing.T) {
		ctx := context.Background()

		ctx = PutUserIDIntoContext(ctx, "123")

		userID, err := GetUserIDFromContext(ctx)
		require.NoError(t, err)
		require.Equal(t, "123", userID)
	})
}
