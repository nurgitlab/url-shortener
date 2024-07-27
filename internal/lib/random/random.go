package random

import (
	"math/rand"
)

func NewRandomString(n int) string {
	ans := ""

	for i := 0; i < n; i++ {
		ans += string(byte(97 + rand.Intn(25)))
	}

	return ans
}
