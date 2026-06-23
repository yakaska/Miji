package slug

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Base62SlugGenerator struct {
	length int
}

func NewBase62SlugGenerator(length int) *Base62SlugGenerator {
	return &Base62SlugGenerator{length: length}
}

func (g *Base62SlugGenerator) Generate() (string, error) {
	result := make([]byte, g.length)
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < g.length; i++ {
		num, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[num.Int64()]
	}

	return string(result), nil
}
