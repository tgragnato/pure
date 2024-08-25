package checks

import "testing"

func TestCheckDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		domain string
		want   bool
	}{
		{
			domain: "apple-finance.query.yahoo.com.",
			want:   true,
		},
		{
			domain: "www.semrush.com.",
			want:   true,
		},
		{
			domain: "dns4torpnlfs2ifuz2s2yf3fc7rdmsbhm6rw75euj35pac6ap25zgqad.onion.",
			want:   false,
		},
		{
			domain: "tgragnato.it.",
			want:   true,
		},
		{
			domain: "test.tgragnato.it.",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			if got := CheckDomain(tt.domain); got != tt.want {
				t.Errorf("CheckDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
