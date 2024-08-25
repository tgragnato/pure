package checks

import (
	"net"
	"testing"
)

func TestGeoChecks_CheckIPs(t *testing.T) {
	t.Parallel()

	g := &GeoChecks{}
	g.initCidr()

	tests := []struct {
		name string
		ips  []net.IP
		want bool
	}{
		{
			name: "Test NULL IPv4",
			ips:  []net.IP{net.ParseIP("0.0.0.0")},
			want: false,
		},
		{
			name: "Test NULL IPv6",
			ips:  []net.IP{net.ParseIP("::")},
			want: false,
		},
		{
			name: "Test loopback IPv4",
			ips:  []net.IP{net.ParseIP("127.0.0.1")},
			want: false,
		},
		{
			name: "Test loopback IPv6",
			ips:  []net.IP{net.ParseIP("::1")},
			want: false,
		},
		{
			name: "Test private IPv4",
			ips:  []net.IP{net.ParseIP("192.168.1.255")},
			want: false,
		},
		{
			name: "Test private IPv6",
			ips:  []net.IP{net.ParseIP("fc00::1")},
			want: false,
		},
		{
			name: "Test public IPv4",
			ips:  []net.IP{net.ParseIP("1.1.1.1")},
			want: true,
		},
		{
			name: "Test public IPv6",
			ips:  []net.IP{net.ParseIP("2606:4700:4700::1111")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := g.CheckIPs(tt.ips); got != tt.want {
				t.Errorf("GeoChecks.CheckIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}
