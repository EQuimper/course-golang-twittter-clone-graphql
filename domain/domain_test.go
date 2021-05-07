package domain

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestMain(t *testing.M) {
	passwordCost = bcrypt.MinCost

	os.Exit(t.Run())
}
