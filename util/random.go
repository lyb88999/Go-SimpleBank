package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

}

// RandInt generate a random integer between min and max
func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 -> max-min
}

// RandString generate a random string of length n
func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandOwner generate a random owner name
func RandOwner() string {
	return RandString(6)
}

// RandMoney generate a random money amount
func RandMoney() int64 {
	return RandInt(0, 1000)
}

// RandomCurrency generate a random currency
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
