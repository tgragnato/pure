package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	proxy, _   = url.Parse("socks5://[::1]:9050")
	httpClient = &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   time.Second,
			KeepAlive: time.Minute,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          4,
		MaxIdleConnsPerHost:   4,
		MaxConnsPerHost:       4,
		IdleConnTimeout:       time.Minute,
		TLSHandshakeTimeout:   time.Second,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: time.Second,
		DisableKeepAlives:     false,
	}}
	urls = []string{
		"https://egsro.com",
		"https://egsro.com/_fd",
		"https://egsro.com/_fd/",
		"https://egsro.com/_tr",
		"https://egsro.com/_tr/",
		"https://egsro.com/_zc",
		"https://egsro.com/_zc/",
		"https://egsro.com/",
		"https://egsro.com/brKcJyirt.js",
		"https://egsro.com/favicon.ico",
		"https://egsro.com/robots.txt",
		"https://tinyurl.com/HIMALA12/",
		"https://tinyurl.com/htsoni/",
		"https://tinyurl.com/Loura41r/",
	}
	domains = []string{
		"http://694east.com",
		"http://abortionclinician.com",
		"http://adanadaemlakci.com",
		"http://alohateriyakigrillct.com",
		"http://amirlabd.com",
		"http://andyalcoholic.com",
		"http://apedxroll.com",
		"http://apexroll.com",
		"http://aqounoli.com",
		"http://arthurnascosh.com",
		"http://av82com.com",
		"http://bantosas.com",
		"http://barbaraton.com",
		"http://biken.site",
		"http://brikonagio.net",
		"http://cbdido.com",
		"http://clmm1.com",
		"http://confortour.com",
		"http://cukuinhospital.com",
		"http://diva138.com",
		"http://dkholpok.com",
		"http://doorstoda.com",
		"http://einfachsuchen.com",
		"http://escapepodnet.com",
		"http://fotografwest.com",
		"http://francegiude.com",
		"http://handalh.com",
		"http://ilyion.com",
		"http://imtoken-app.com",
		"http://irishwoodwatches.com",
		"http://jzcai1603.com",
		"http://kanboris.com",
		"http://kasvigas.com",
		"http://kasvigas.net",
		"http://kolgyh.com",
		"http://labuiand.com",
		"http://lagartad.com",
		"http://lakartpa.com",
		"http://larabonis.com",
		"http://latfitx.com",
		"http://linioki.com",
		"http://lintosbant.com",
		"http://madishwm.com",
		"http://maikodmi.net",
		"http://malawais.com",
		"http://malisdakoni.com",
		"http://mandarroskolo.com",
		"http://markiisa.com",
		"http://masandinara.com",
		"http://maxicomsrl.com",
		"http://mazaratnadik.com",
		"http://mzolpon.com",
		"http://mzyikon.com",
		"http://oy2order.com",
		"http://parteekit.com",
		"http://point-s.net",
		"http://pucuk69.com",
		"http://servantidal.com",
		"http://startkolcas.com",
		"http://trkoli.com",
		"http://tuv-nord.net",
		"http://viohki.com",
		"http://viokyi.com",
		"http://waletnas.com",
		"http://zarbonk.com",
		"http://zoheoki.com",
	}
	paths = []string{
		"",
		"/",
		"/4HXjoo80791dYtM1195zuwviakgln8171EJFZXXVQGNPXAOK13839/18625R22#pdbppwwlaexditvthghaj",
		"/4Iotjw81001AlPk1363ebsqyzwhvs10521SZXAOEISQDZMGFQ2040632/18649v11#qfhofadsqooaihokwkokk",
		"/4sCtWa80791GpEG1195teuwrvsxbm8171MYTFABBBOQCNHXM13839/18625v22#riwkuhwealeouqzihzvcb",
		"/4uIEqG80984Plej1363mrquqxbpan10521NFBEGXOICSJWGLT2040632/18649e11#uhrmkzosemksdrbqpzucd",
		"/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22#drsbyizemphcffpuxqwje",
		"/4WKKpm80770NJNV1619ryqbstiovr8171MWJCRLTEIYFSZON13839/18625R22#lhsglrclhqluewhayihdy",
		"/5dWwxC80791dfpi1195dkmhnmcxcg8171ASVGXUTOSQPQBUC13839/18625c22#nfocppzyzkhclaigfdedr",
		"/5eYbZo80702JGjN767ypgmbwniii8171AWRGFKLKLRPKLPM13839/18625X22#gypdyeghnejaegotzgqyt",
		"/5fJPcu80770LrUx1619qorpaodvmu8171ALLOFCDQMYEKMIT13839/18625x22#xqxhaflofarvykfyvweyr",
		"/5laNXF80790UzLG767fwqmqfwwsz8171JFPBMGWFWHOFAWE13839/18625r22#fbtzrkfnduqlqqiegsvcg",
		"/5NvKDQ80984riYX1363umvrmnphby10521FORIIULLATURFLG2040632/18649Z11#pfxglimjeauhwqekshuix",
		"/5TcdFH80792pQrq1362smjculxvav8171LYHBVZZVECVIUNG13839/18625M22#zuylgqjyjiflobmwjfcvk",
		"/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"/6DUEvB80790jvRD767claopcblok8171BPLORCVMSAKIGBA13839/18625i22",
		"/6jxIcg80770jiBE1619zosbmbjiag8171IJNXYRAQAPYTMBB13839/18625T22",
		"/6RBIAx80791acSe1195cfqticzrpx8171RUIUVYDAHYVPUIZ13839/18625n22",
		"/6UPeBZ80984EQhU1363qhpdnzbsew10521NBJBRWGNILCFBEZ2040632/18649i11",
		"/6yfdln80792fpEQ1362mzednzrstj8171YKZXKOVZXAHYVAE13839/18625N22#lqzxddadypshbybofumzp",
		"/about",
		"/news",
		"/opt-out",
		"/privacy",
		"/robots.txt",
		"/t/robots.txt",
		"/track/3TMvbV80702vdqw767brbxkphbsk8171SFFDKUUQMREBQPG13839/18625p22",
	}
	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/128.0",
	}
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for range time.NewTicker(time.Minute).C {
		items := make([]string, len(urls))
		copy(items, urls)
		for _, domain := range domains {
			for _, path := range paths {
				items = append(items, domain+path)
			}
		}
		rand.Shuffle(len(urls), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		for index, url := range items {
			go call(url, userAgents[index%len(userAgents)])
		}
	}

	<-signalCh
}

func call(url string, userAgent string) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
