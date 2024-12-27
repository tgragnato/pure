package spam

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	proxy, _   = url.Parse("socks5://[::1]:9050")
	httpClient = &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   time.Minute,
			KeepAlive: time.Minute,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          250,
		MaxIdleConnsPerHost:   5,
		MaxConnsPerHost:       5,
		IdleConnTimeout:       time.Minute,
		TLSHandshakeTimeout:   time.Minute,
		ExpectContinueTimeout: time.Minute,
		ResponseHeaderTimeout: time.Minute,
		DisableKeepAlives:     false,
	}}
	proxiedPrefixes = []string{
		"http://142.202.189.213",
		"http://162.62.134.81",
		"http://162.62.230.106",
		"http://178.255.74.137",
		"http://2.59.119.3",
		"http://chetariffa.it",
		"https://5u9c.r.ag.d.sendibm3.com",
		"https://acr.emails.zephiromedia.it",
		"https://affiliate.across.it",
		"https://blog.across.it",
		"https://click.e.dyson.it",
		"https://coreg.across.it",
		"https://dashboard.across.it",
		"https://dbunico.across.it",
		"https://ftps.across.it",
		"https://lead.across.it",
		"https://lead2.across.it",
		"https://mooo.com",
		"https://network.across.it",
		"https://pacc.across.it",
		"https://plan.across.it",
		"https://privacy.dyson.com/it",
		"https://shortener.across.it",
		"https://view.e.dyson.it",
		"https://www.across.it",
		"https://www.dyson.it/unsubscribe-from-dyson-emails",
		"https://www.ediscom.it",
		"https://www.promo-home.com",
		"https://www.promo-home.com/it/climatizzatore-ariel-detrazioni-omaggio-ms",
		"https://zephiromedia.it",
	}
	paths = []string{
		"",
		"/",
		"/about",
		"/favicon.ico",
		"/informativa-policy/",
		"/mail/contact_me.php",
		"/news",
		"/opt-out",
		"/privacy-be-nl",
		"/privacy-de",
		"/privacy-es",
		"/privacy-pl",
		"/privacy-policy-ita",
		"/privacy-se",
		"/privacy",
		"/robots.txt",
		"/t/robots.txt",
		"/tr.php",
	}
	userAgents = []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.6778.69 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.6778.69 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.6778.69 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.6778.69 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; rv:133.0) Gecko/20100101 Firefox/133.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/133.0",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/133.0.0.0 Safari/537.36",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/133.0.6778.69 Safari/537.36",
	}
	Counter = uint64(0)
)
