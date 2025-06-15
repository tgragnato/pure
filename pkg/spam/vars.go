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
		"http://m.bib.us.es",
		"https://apis.be-mail.it",
		"https://cdn.files-text.com/us-south1/api/lc/att/1520",
		"https://email.analystratings.net",
		"https://s3-us-west-1.amazonaws.com/gbso73q0n9vvsozerj9uhk4i5o0g46",
		"https://sms.be-mail.it",
		"https://www.ketchupadv.com",
		//"http://142.202.189.213",
		//"http://162.62.134.81",
		//"http://162.62.230.106",
		//"http://178.255.74.137",
		//"http://2.59.119.3",
	}
	paths = []string{
		"",
		"/",
		"/4webc1ohqnif08kbbue9688u8c8nuy/0qksk006gazecddnvepmn58vn7npty.html",
		"/4webc1ohqnif08kbbue9688u8c8nuy/0qksk00dgazecddnvesmn58vn7njty.html",
		"/4webc1ohqnif08kbbue9688u8c8nuy/701jjle7q1cpmpdv9516xz8z08wjqo.jpg",
		"/4webc1ohqnif08kbbue9688u8c8nuy/701ljle7q1cpmpdv9506xz8z08wjqo.jpg",
		"/a4deb57133efd03363f2687b34f66df6/102276%20EN.png",
		"/a4deb57133efd03363f2687b34f66df6/102276%20ES.png",
		"/a4deb57133efd03363f2687b34f66df6/102276%20IT.png",
		"/about",
		"/favicon.ico",
		"/informativa-policy/",
		"/ita",
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
		"AppleCoreMedia/1.0.0.23A344 (Macintosh; U; Intel Mac OS X 14_0; da_dk)",
		"Dalvik/2.1.0 (Linux; U; Android 11; Tibuta_MasterPad-E100 Build/RP1A.201005.006)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 15_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/19G82 Instagram 306.0.0.20.118 (iPhone12,1; iOS 15_6_1; en_GB; en; scale=2.00; 828x1792; 529083166) NW/3",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [LinkedInApp]/9.",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 [LinkedInApp]/9.28.7586",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 18_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 11; SM-A115M Build/RP1A.200720.012; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/102.0.5005.125 Mobile Safari/537.36 Instagram 306.0.0.35.109 Android (30/11; 280dpi; 720x1411; samsung; SM-A115M; a11q; qcom; pt_BR; 530130405)",
		"Mozilla/5.0 (Linux; Android 13; SAMSUNG SM-T220) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/23.0 Chrome/115.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 13; SM-F711U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36 EdgA/114.0.1823.43",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G532MT Build/MMB29T; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/99.0.4844.88 Mobile Safari/537.36 [FB_IAB/FB4A;FBAV/436.0.0.35.101;]",
		"Mozilla/5.0 (Linux; Android 9) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/119.0.6045.66 Mobile DuckDuckGo/1 Lilo/1.2.3 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; U; PPC Mac OS X Mach-O; en-US; rv:1.7.6) Gecko/20050319",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/128.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76 GLS/97.10.7399.100",
		"Mozilla/5.0 (X11; Linux x86_64; SMARTEMB Build/3.12.9076) AppleWebKit/537.36 (KHTML, like Gecko) Chromium/103.0.5060.129 Chrome/103.0.5060.129 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/102.0.5143.178 Chrome/102.0.5143.178 Safari/537.36",
		"Mozilla/5.0 (X11; U; Linux i586; en-US; rv:1.0.0) Gecko/20020623 Debian/1.0.0-0.woody.1",
		"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.2.1) Gecko/20021208 Debian/1.2.1-2",
		"Mozilla/5.0 (X11; U; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/115.0.5738.217 Chrome/115.0.5738.217 Safari/537.36",
	}
	Counter = uint64(0)
)
