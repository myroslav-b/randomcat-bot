package main

import "testing"

func TestMatchTocken(t *testing.T) {
	cases := []struct {
		token string
		want  bool
	}{
		{"1223166507:AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", true},
		{"1223166500:-AFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6_", true},
		{"1223166507AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", false},
		{"223166507:AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", false},
		{"1223166507:AAFHcxGFvCC_zVh6r-XZ_bj6CchwXMUw6k", false},
		{":AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", false},
		{"1223166507:", false},
		{"AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", false},
		{"1223166507", false},
		{"", false},
		{"1223166507:AAFHcxGFvCC+zVh6r-XZ_bj6CchwXM-Uw6k", false},
		{"1223A66507:AAFHcxGFvCC_zVh6r-XZ_bj6CchwXM-Uw6k", false},
	}
	for _, c := range cases {
		got := matchToken(c.token)
		if got != c.want {
			t.Errorf("matchToken(%q) == %v, want %v", c.token, got, c.want)
		}
	}
}
