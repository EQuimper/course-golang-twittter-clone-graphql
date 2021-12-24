package faker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/equimper/twitter/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Password equal password
var Password = "$2a$04$5vcLELOMsTLvXsUM9Nd.FegPIrNRT3s1mjEyH3rk2mh/b9iQGHLzG"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func randStringLowerRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}

	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Username() string {
	return randStringRunes(RandInt(2, 10))
}

func ID() string {
	return fmt.Sprintf("%s-%s-%s-%s", randStringRunes(4), randStringRunes(4), randStringRunes(4), randStringRunes(4))
}

func UUID() string {
	return uuid.Generate()
}

func Email() string {
	return fmt.Sprintf("%s@example.com", randStringLowerRunes(RandInt(5, 10)))
}

func RandStr(n int) string {
	return randStringRunes(n)
}
