package main

import "testing"

func TestMatchPort(t *testing.T) {
	cases := []struct {
		port string
		want bool
	}{
		{"1", true},
		{"12345", true},
		{"1234", true},
		{"65534", true},
		{"65536", false},
		{"0", false},
		{"0123", false},
		{"123456", false},
		{"55a55", false},
		{"", false},
	}
	for _, c := range cases {
		got := matchPort(c.port)
		if got != c.want {
			t.Errorf("matchToken(%q) == %v, want %v", c.port, got, c.want)
		}
	}
}
