package service

import "testing"

func TestBase62(t *testing.T) {
	tests := []struct {
		name   string
		n      uint64
		length int
		want   string
	}{
		{name: "zero is padded", n: 0, length: 4, want: "aaaa"},
		{name: "one is padded", n: 1, length: 4, want: "aaab"},
		{name: "base boundary", n: 62, length: 4, want: "aaba"},
		{name: "maximum fitting value", n: 62*62 - 1, length: 2, want: "99"},
		{name: "overflow", n: 62 * 62, length: 2, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base62(tt.n, tt.length); got != tt.want {
				t.Fatalf("Base62(%d, %d) = %q, want %q", tt.n, tt.length, got, tt.want)
			}
		})
	}
}

func TestShortenURL(t *testing.T) {
	if _, err := ShortenURL(-1); err == nil {
		t.Fatal("ShortenURL(-1) returned no error")
	}

	hash, err := ShortenURL(42)
	if err != nil {
		t.Fatalf("ShortenURL(42) returned an error: %v", err)
	}
	if len(hash) != 8 {
		t.Fatalf("hash length = %d, want 8", len(hash))
	}
	if got, want := hash[4:], "aaaQ"; got != want {
		t.Fatalf("encoded ID = %q, want %q", got, want)
	}
}
