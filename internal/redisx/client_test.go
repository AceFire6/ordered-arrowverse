package redisx

import "testing"

func TestRedactPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "redis url with password",
			in:   "redis://:secret@redis.example.com:6379/0",
			want: "redis://@redis.example.com:6379/0",
		},
		{
			name: "redis url with user:password",
			in:   "redis://user:secret@redis.example.com:6379/0",
			want: "redis://user@redis.example.com:6379/0",
		},
		{
			name: "redis url without password",
			in:   "redis://redis.example.com:6379/0",
			want: "redis://redis.example.com:6379/0",
		},
		{
			name: "non-redis url also gets scrubbed",
			in:   "postgres://user:secret@host/db",
			want: "postgres://user@host/db",
		},
		{
			name: "empty",
			in:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactPassword(tt.in); got != tt.want {
				t.Errorf("redactPassword(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
