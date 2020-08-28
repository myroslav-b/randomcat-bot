package main

import "testing"

func TestMatchIP(t *testing.T) {
	cases := []struct {
		ip   string
		want bool
	}{
		{"192.168.0.1", true},
		{"127.0.0.1", true},
		{"10.1.1.1", true},
		{"1.10.100.0", true},
		{"192.168.0.1.", false},
		{"222.192.168.0.1", false},
		{"192.1680.0.1", false},
		{"192.168.257.1", false},
		{"192.168.00.1", false},
		{"192. 168.0.1", false},
		{"192.168..1", false},
		{"192.168.1", false},
		{"192-168-1-1", false},
		{"1921681", false},
		{"", false},
	}
	for _, c := range cases {
		got := matchIP(c.ip)
		if got != c.want {
			t.Errorf("matchIP(%q) == %v, want %v", c.ip, got, c.want)
		}
	}
}
