package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate a radom integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[(rand.Intn(k))]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a radom owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomMoney generate a radmom currency code
func RadomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random eamil
func RandomEmail() string {
	return fmt.Sprintf("%sgmail.com", RandomString(6))
}

// One to many: which basically means: 1 user can have multiple accounts, but 1 account can only belong to exactly 1 single user
