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
		"https://acrosau.be-mail.it",
		"https://adfi.be-mail.it",
		"https://ads.ketchupadv.com",
		"https://adv-console.be-mail.it",
		"https://affiliate.across.it",
		"https://apis.be-mail.it",
		"https://assets.be-mail.it",
		"https://bci.be-mail.it",
		"https://bid.ketchupadv.com",
		"https://blog.across.it",
		"https://budget.ketchupadv.com",
		"https://click.e.dyson.it",
		"https://coreg.across.it",
		"https://dashboard.across.it",
		"https://db.ketchupadv.com",
		"https://dbunico.across.it",
		"https://dev.api.planner.ketchupadv.com",
		"https://dev.budget.ketchupadv.com",
		"https://dev.console.be-mail.it",
		"https://dev.ketchupadv.com",
		"https://dev.nobounce.be-mail.it",
		"https://dev.planner.ketchupadv.com",
		"https://dev.sms.be-mail.it",
		"https://flowly.ketchupadv.com",
		"https://ftps.across.it",
		"https://ingroferrendpoint.ketchupadv.com",
		"https://lamaisondelamour.ketchupadv.com",
		"https://lead.across.it",
		"https://lead2.across.it",
		"https://monitor.ketchupadv.com",
		"https://mooo.com",
		"https://multiquadri.ketchupadv.com",
		"https://native-console.be-mail.it",
		"https://network.across.it",
		"https://news.ketchupadv.com",
		"https://nobounce.be-mail.it",
		"https://one-eventi.ketchupadv.com",
		"https://pacc.across.it",
		"https://plan.across.it",
		"https://pma.ketchupadv.com",
		"https://privacy.dyson.com/it",
		"https://secure.ketchupadv.com",
		"https://shortener.across.it",
		"https://sms.be-mail.it",
		"https://svn.ketchupadv.com",
		"https://view.e.dyson.it",
		"https://www.across.it",
		"https://www.deepseek.com",
		"https://www.dyson.it/unsubscribe-from-dyson-emails",
		"https://www.ediscom.it",
		"https://www.ketchupadv.com",
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
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.204 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.204 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.6778.204 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	}
	Counter = uint64(0)
)
