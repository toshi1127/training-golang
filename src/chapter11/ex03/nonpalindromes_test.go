package palindrom

import (
	"math/rand"
	"testing"
	"time"
)

func randomNonPalindrome(rng *rand.Rand) string {
	n := 2 + rng.Intn(23) // 24までのランダムな長さ
	runes := make([]rune, n)
	for i := 0; i < n-1; i++ {
		r := rune(rng.Intn(0x0999)) // '\u0998' までのランダムなルーン
		runes[i] = r
	}
	runes[n-1] = '\u0999'
	return string(runes)
}

func TestRandomNonPalindrome(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomNonPalindrome(rng)
		if IsPalindrome(p) {
			t.Errorf("IsPandrome(%q) = true", p)
		}
	}
}