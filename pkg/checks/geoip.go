package checks

import (
	"net"

	"gitlab.torproject.org/tpo/anti-censorship/geoip"
)

type GeoChecks struct {
	geo          *geoip.Geoip
	unrouteables [28]*net.IPNet
}

func (g *GeoChecks) CheckIPs(ips []net.IP) bool {
	for x := range ips {
		if !ips[x].IsGlobalUnicast() {
			return false
		}
		for y := range g.unrouteables {
			if g.unrouteables[y].Contains(ips[x]) {
				return false
			}
		}

		if g.geo != nil {
			country, ok := g.geo.GetCountryByAddr(ips[x])
			if !ok {
				return false
			}
			switch country {
			case "CN", "HK", "MO", "RU", "BY", "IR", "KP":
				return false
			}
		}
	}

	return true
}

func (g *GeoChecks) GetCountry(ip net.IP) string {
	if country, ok := g.geo.GetCountryByAddr(ip); ok {
		return country
	}
	return ""
}
