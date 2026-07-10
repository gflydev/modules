package storage

import "testing"

func TestIsSafeRelPath(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"avatars", true},
		{"avatars/2024", true},
		{"a/b/c", true},
		{"", false},
		{"..", false},
		{"../etc", false},
		{"avatars/../../etc", false},
		{"/etc/passwd", false},
		{"/abs", false},
		{"a/../../b", false},
		{"with\x00null", false},
	}

	for _, c := range cases {
		if got := IsSafeRelPath(c.in); got != c.want {
			t.Errorf("IsSafeRelPath(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestIsSafeFileName(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"avatar.png", true},
		{"my-file_01.jpg", true},
		{"", false},
		{"../avatar.png", false},
		{"dir/avatar.png", false},
		{"dir\\avatar.png", false},
		{"..", false},
	}

	for _, c := range cases {
		if got := IsSafeFileName(c.in); got != c.want {
			t.Errorf("IsSafeFileName(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
