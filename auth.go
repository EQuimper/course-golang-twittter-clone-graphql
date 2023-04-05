package twitter

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 2
	PasswordMinLength = 6
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
	Login(ctx context.Context, input LoginInput) (AuthResponse, error)
}

type AuthTokenService interface {
	CreateAccessToken(ctx context.Context, user User) (string, error)
	CreateRefreshToken(ctx context.Context, user User, tokenID string) (string, error)
	ParseToken(ctx context.Context, payload string) (AuthToken, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
}

type AuthToken struct {
	ID  string
	Sub string
}

type AuthResponse struct {
	AccessToken string
	User        User
}

type RegisterInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

func (in *RegisterInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)

	in.Username = strings.TrimSpace(in.Username)
}

func (in RegisterInput) Validate() error {
	if len(in.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username not long enough, (%d) characters at least", ErrValidation, UsernameMinLength)
	}

	if !emailRegexp.MatchString(in.Email) {
		return fmt.Errorf("%w: email not valid", ErrValidation)
	}

	if len(in.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password not long enough, (%d) characters at least", ErrValidation, PasswordMinLength)
	}

	if in.Password != in.ConfirmPassword {
		return fmt.Errorf("%w: confirm password must match the password", ErrValidation)
	}

	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (in *LoginInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)
}

func (in LoginInput) Validate() error {
	if !emailRegexp.MatchString(in.Email) {
		return fmt.Errorf("%w: email not valid", ErrValidation)
	}

	if len(in.Password) < 1 {
		return fmt.Errorf("%w: password required", ErrValidation)
	}

	return nil
}
