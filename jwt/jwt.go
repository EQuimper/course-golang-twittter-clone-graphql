package jwt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/equimper/twitter"
	"github.com/equimper/twitter/config"
	"github.com/lestrrat-go/jwx/jwa"
	jwtGo "github.com/lestrrat-go/jwx/jwt"
)

var signatureType = jwa.HS256

var now = time.Now

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (twitter.AuthToken, error) {
	token, err := jwtGo.ParseRequest(
		r,
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)
	if err != nil {
		return twitter.AuthToken{}, twitter.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func buildToken(token jwtGo.Token) twitter.AuthToken {
	return twitter.AuthToken{
		ID:  token.JwtID(),
		Sub: token.Subject(),
	}
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (twitter.AuthToken, error) {
	token, err := jwtGo.Parse(
		[]byte(payload),
		jwtGo.WithValidate(true),
		jwtGo.WithIssuer(ts.Conf.JWT.Issuer),
		jwtGo.WithVerify(signatureType, []byte(ts.Conf.JWT.Secret)),
	)
	if err != nil {
		return twitter.AuthToken{}, twitter.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user twitter.User, tokenID string) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, twitter.RefreshTokenLifetime, ts.Conf); err != nil {
		return "", err
	}

	if err := t.Set(jwtGo.JwtIDKey, tokenID); err != nil {
		return "", fmt.Errorf("error set jwt id: %v", err)
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign jwt: %v", err)
	}

	return string(token), nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user twitter.User) (string, error) {
	t := jwtGo.New()

	if err := setDefaultToken(t, user, twitter.AccessTokenLifetime, ts.Conf); err != nil {
		return "", err
	}

	token, err := jwtGo.Sign(t, signatureType, []byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign jwt: %v", err)
	}

	return string(token), nil
}

func setDefaultToken(t jwtGo.Token, user twitter.User, lifetime time.Duration, conf *config.Config) error {
	if err := t.Set(jwtGo.SubjectKey, user.ID); err != nil {
		return fmt.Errorf("error set jwt sub: %v", err)
	}

	if err := t.Set(jwtGo.IssuerKey, conf.JWT.Issuer); err != nil {
		return fmt.Errorf("error set jwt issuer key: %v", err)
	}

	if err := t.Set(jwtGo.IssuedAtKey, now().Unix()); err != nil {
		return fmt.Errorf("error set jwt issued at key: %v", err)
	}

	if err := t.Set(jwtGo.ExpirationKey, now().Add(lifetime).Unix()); err != nil {
		return fmt.Errorf("error set jwt expired at: %v", err)
	}

	return nil
}
