package sni

import "testing"

func Test_getHostPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		sni  string
		want string
	}{
		{
			sni:  "p27-imap.mail.me.com",
			want: "p27-imap.mail.me.com:993",
		},
		{
			sni:  "imap.gmail.com",
			want: "imap.gmail.com:993",
		},
		{
			sni:  "example.com",
			want: "example.com:443",
		},
	}
	for _, tt := range tests {
		t.Run(tt.sni, func(t *testing.T) {
			if got := getHostPort(tt.sni); got != tt.want {
				t.Errorf("getHostPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
