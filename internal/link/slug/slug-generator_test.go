package slug_test

import (
	"Miji/internal/link/slug"
	"strings"
	"testing"
)

func TestBase62SlugGenerator_Generate(t *testing.T) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	tests := []struct {
		name   string
		length int
	}{
		{"short", 6},
		{"standard", 8},
		{"long", 16},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := slug.NewBase62SlugGenerator(tc.length)

			s, err := g.Generate()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(s) != tc.length {
				t.Errorf("expected length %d, got %d", tc.length, len(s))
			}
			for _, ch := range s {
				if !strings.ContainsRune(alphabet, ch) {
					t.Errorf("character %q not in alphabet", ch)
				}
			}
		})
	}
}
