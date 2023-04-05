package jwt

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/config"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
)

var (
	conf         *config.Config
	tokenService *TokenService
)

func TestMain(m *testing.M) {
	config.LoadEnv(".env.test")

	conf = config.New()

	tokenService = NewTokenService(conf)

	os.Exit(m.Run())
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("can create a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithVerify(signatureType, []byte(conf.JWT.Secret)),
			jwtGo.WithIssuer(conf.JWT.Issuer),
		)
		require.NoError(t, err)

		require.Equal(t, "123", tok.Subject())
		require.Equal(t, now().Add(twitter.AccessTokenLifetime).Unix(), tok.Expiration().Unix())

		teardownTimeNow(t)
	})
}

func TestTokenService_CreateRefreshToken(t *testing.T) {
	t.Run("can create a valid refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateRefreshToken(ctx, user, "456")
		require.NoError(t, err)

		now = func() time.Time {
			return time.Now()
		}

		tok, err := jwtGo.Parse(
			[]byte(token),
			jwtGo.WithValidate(true),
			jwtGo.WithVerify(signatureType, []byte(conf.JWT.Secret)),
			jwtGo.WithIssuer(conf.JWT.Issuer),
		)
		require.NoError(t, err)

		require.Equal(t, "123", tok.Subject())
		require.Equal(t, "456", tok.JwtID())
		require.Equal(t, now().Add(twitter.RefreshTokenLifetime).Unix(), tok.Expiration().Unix())

		teardownTimeNow(t)
	})
}

func TestTokenService_ParseToken(t *testing.T) {
	t.Run("can parse a valid access token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, "123", tok.Sub)
	})

	t.Run("can parse a valid refresh token", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		token, err := tokenService.CreateRefreshToken(ctx, user, "456")
		require.NoError(t, err)

		tok, err := tokenService.ParseToken(ctx, token)
		require.NoError(t, err)
		require.Equal(t, "123", tok.Sub)
		require.Equal(t, "456", tok.ID)
	})

	t.Run("return err if invalid access token", func(t *testing.T) {
		ctx := context.Background()

		_, err := tokenService.ParseToken(ctx, "invalid token")
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)
	})

	t.Run("return err if access token expired", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		now = func() time.Time {
			return time.Now().Add(-twitter.AccessTokenLifetime * 5)
		}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		_, err = tokenService.ParseToken(ctx, token)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)

		teardownTimeNow(t)
	})
}

func TestTokenService_ParseTokenFromRequest(t *testing.T) {
	t.Run("can parse an access token from the request", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		req := httptest.NewRequest("GET", "/", nil)

		accessToken, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		req.Header.Set("Authorization", accessToken)

		token, err := tokenService.ParseTokenFromRequest(ctx, req)
		require.NoError(t, err)

		require.Equal(t, "123", token.Sub)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		token, err = tokenService.ParseTokenFromRequest(ctx, req)
		require.NoError(t, err)

		require.Equal(t, "123", token.Sub)
	})

	t.Run("expired token will fail", func(t *testing.T) {
		ctx := context.Background()
		user := twitter.User{
			ID: "123",
		}

		req := httptest.NewRequest("GET", "/", nil)

		now = func() time.Time {
			return time.Now().Add(-twitter.AccessTokenLifetime * 5)
		}

		accessToken, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)

		req.Header.Set("Authorization", accessToken)

		_, err = tokenService.ParseTokenFromRequest(ctx, req)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)

		teardownTimeNow(t)
	})

	t.Run("expired token will fail", func(t *testing.T) {
		ctx := context.Background()

		req := httptest.NewRequest("GET", "/", nil)

		req.Header.Set("Authorization", "invalid token")

		_, err := tokenService.ParseTokenFromRequest(ctx, req)
		require.ErrorIs(t, err, twitter.ErrInvalidAccessToken)
	})
}

func teardownTimeNow(t *testing.T) {
	t.Helper()

	now = func() time.Time {
		return time.Now()
	}
}
