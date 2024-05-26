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
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       time.Hour,
		TLSHandshakeTimeout:   time.Second,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: time.Second,
		DisableKeepAlives:     false,
	}}
	numberOfWorkers = 10
	urls            = []string{
		"https://tinyurl.com/HIMALA12/4GbkAE80702KUYR767luahfovgfo8171TSOJJDIWIOPSUFW13839/18625t22#txzrkskcjttglrmclxwfu",
		"http://linioki.com/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"https://www.httpsimage.com/v2/03d8266e-9bcd-49b3-b11f-57360eafbe31.png",
		"https://tinyurl.com/HIMALA12/5eYbZo80702JGjN767ypgmbwniii8171AWRGFKLKLRPKLPM13839/18625X22#gypdyeghnejaegotzgqyt",
		"http://linioki.com/5eYbZo80702JGjN767ypgmbwniii8171AWRGFKLKLRPKLPM13839/18625X22#gypdyeghnejaegotzgqyt",
		"https://tinyurl.com/HIMALA12/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"http://linioki.com/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"https://tinyurl.com/HIMALA12/track/3TMvbV80702vdqw767brbxkphbsk8171SFFDKUUQMREBQPG13839/18625p22",
		"http://linioki.com/track/3TMvbV80702vdqw767brbxkphbsk8171SFFDKUUQMREBQPG13839/18625p22",
		"https://tinyurl.com/HIMALA12/4WKKpm80770NJNV1619ryqbstiovr8171MWJCRLTEIYFSZON13839/18625R22#lhsglrclhqluewhayihdy",
		"http://linioki.com/4WKKpm80770NJNV1619ryqbstiovr8171MWJCRLTEIYFSZON13839/18625R22#lhsglrclhqluewhayihdy",
		"https://tinyurl.com/HIMALA12/4awNPL80770yTrT1619tgnqwmdjpt8171AAXOREBDUDUWDAA13839/18625g22#nqrnyfewonnhvgrizlsjz",
		"http://linioki.com/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22#drsbyizemphcffpuxqwje",
		"https://tinyurl.com/HIMALA12/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22#drsbyizemphcffpuxqwje",
		"http://linioki.com/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22#drsbyizemphcffpuxqwje",
		"https://tinyurl.com/HIMALA12/5fJPcu80770LrUx1619qorpaodvmu8171ALLOFCDQMYEKMIT13839/18625x22#xqxhaflofarvykfyvweyr",
		"http://linioki.com/5fJPcu80770LrUx1619qorpaodvmu8171ALLOFCDQMYEKMIT13839/18625x22#xqxhaflofarvykfyvweyr",
		"https://tinyurl.com/HIMALA12/6jxIcg80770jiBE1619zosbmbjiag8171IJNXYRAQAPYTMBB13839/18625T22",
		"http://linioki.com/6jxIcg80770jiBE1619zosbmbjiag8171IJNXYRAQAPYTMBB13839/18625T22",
		"https://tinyurl.com/HIMALA12/4HXjoo80791dYtM1195zuwviakgln8171EJFZXXVQGNPXAOK13839/18625R22#pdbppwwlaexditvthghaj",
		"http://linioki.com/4HXjoo80791dYtM1195zuwviakgln8171EJFZXXVQGNPXAOK13839/18625R22#pdbppwwlaexditvthghaj",
		"https://tinyurl.com/HIMALA12/4sCtWa80791GpEG1195teuwrvsxbm8171MYTFABBBOQCNHXM13839/18625v22#riwkuhwealeouqzihzvcb",
		"http://linioki.com/4sCtWa80791GpEG1195teuwrvsxbm8171MYTFABBBOQCNHXM13839/18625v22#riwkuhwealeouqzihzvcb",
		"https://tinyurl.com/HIMALA12/5dWwxC80791dfpi1195dkmhnmcxcg8171ASVGXUTOSQPQBUC13839/18625c22#nfocppzyzkhclaigfdedr",
		"http://linioki.com/5dWwxC80791dfpi1195dkmhnmcxcg8171ASVGXUTOSQPQBUC13839/18625c22#nfocppzyzkhclaigfdedr",
		"https://tinyurl.com/HIMALA12/6RBIAx80791acSe1195cfqticzrpx8171RUIUVYDAHYVPUIZ13839/18625n22",
		"http://linioki.com/6RBIAx80791acSe1195cfqticzrpx8171RUIUVYDAHYVPUIZ13839/18625n22",
		"https://tinyurl.com/htsoni/5TcdFH80792pQrq1362smjculxvav8171LYHBVZZVECVIUNG13839/18625M22#zuylgqjyjiflobmwjfcvk",
		"http://mazaratnadik.com/5TcdFH80792pQrq1362smjculxvav8171LYHBVZZVECVIUNG13839/18625M22#zuylgqjyjiflobmwjfcvk",
		"https://tinyurl.com/htsoni/6yfdln80792fpEQ1362mzednzrstj8171YKZXKOVZXAHYVAE13839/18625N22#lqzxddadypshbybofumzp",
		"http://mazaratnadik.com/6yfdln80792fpEQ1362mzednzrstj8171YKZXKOVZXAHYVAE13839/18625N22#lqzxddadypshbybofumzp",
		"https://tinyurl.com/HIMALA12/5laNXF80790UzLG767fwqmqfwwsz8171JFPBMGWFWHOFAWE13839/18625r22#fbtzrkfnduqlqqiegsvcg",
		"http://linioki.com/5laNXF80790UzLG767fwqmqfwwsz8171JFPBMGWFWHOFAWE13839/18625r22#fbtzrkfnduqlqqiegsvcg",
		"https://tinyurl.com/HIMALA12/6DUEvB80790jvRD767claopcblok8171BPLORCVMSAKIGBA13839/18625i22",
		"http://linioki.com/6DUEvB80790jvRD767claopcblok8171BPLORCVMSAKIGBA13839/18625i22",
		"https://tinyurl.com/Loura41r/4uIEqG80984Plej1363mrquqxbpan10521NFBEGXOICSJWGLT2040632/18649e11#uhrmkzosemksdrbqpzucd",
		"http://madishwm.com/4uIEqG80984Plej1363mrquqxbpan10521NFBEGXOICSJWGLT2040632/18649e11#uhrmkzosemksdrbqpzucd",
		"https://tinyurl.com/Loura41r/5NvKDQ80984riYX1363umvrmnphby10521FORIIULLATURFLG2040632/18649Z11#pfxglimjeauhwqekshuix",
		"http://madishwm.com/5NvKDQ80984riYX1363umvrmnphby10521FORIIULLATURFLG2040632/18649Z11#pfxglimjeauhwqekshuix",
		"https://tinyurl.com/Loura41r/6UPeBZ80984EQhU1363qhpdnzbsew10521NBJBRWGNILCFBEZ2040632/18649i11",
		"http://madishwm.com/6UPeBZ80984EQhU1363qhpdnzbsew10521NBJBRWGNILCFBEZ2040632/18649i11",
		"https://tinyurl.com/Loura41r/4Iotjw81001AlPk1363ebsqyzwhvs10521SZXAOEISQDZMGFQ2040632/18649v11#qfhofadsqooaihokwkokk",
		"http://madishwm.com/4Iotjw81001AlPk1363ebsqyzwhvs10521SZXAOEISQDZMGFQ2040632/18649v11#qfhofadsqooaihokwkokk",
	}
	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	}
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for range time.NewTicker(time.Minute).C {
		for i := 0; i < numberOfWorkers; i++ {
			items := make([]string, len(urls))
			copy(items, urls)
			rand.Shuffle(len(urls), func(i, j int) {
				items[i], items[j] = items[j], items[i]
			})
			for _, url := range items {
				go call(url, i)
			}
		}
	}

	<-signalCh
}

func call(url string, i int) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
