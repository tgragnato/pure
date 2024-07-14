package main

import (
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
		"http://162.62.134.81",
		"http://162.62.230.106",
		"https://jerprint.com",
		"https://panelapi.mooo.com/tr.php",
	}
	directPrefixes = []string{
		"http://162.62.134.81",
		"http://162.62.230.106",
	}
	paths = []string{
		"?token=90962206bffeb5abb89ea2a1bd7b8fad&d=143962&s=820021723",
		"?token=90962206bffeb5abb89ea2a1bd7b8fad&s=820021723&d=143962",
		"?token=bc39323c6c514ab097acb885e06e0713&d=147417&s=851965318",
		"?token=bc39323c6c514ab097acb885e06e0713&s=851965318&d=147417",
		"?token=ecfd39819db7e89b0cbcb7b995282a7a&d=146573&s=840648295",
		"?token=ecfd39819db7e89b0cbcb7b995282a7a&d=146573&s=840648295",
		"",
		"/_fd/",
		"/_tr/",
		"/_zc/",
		"/",
		"/4BVbmS27343JzYU935yrxrokvdvf2280EWOQQAIFKKQKNGH171028/11452c12",
		"/4HXjoo80791dYtM1195zuwviakgln8171EJFZXXVQGNPXAOK13839/18625R22",
		"/4IfNnV82034CQDT833zafgfloxof11723BMQLFBXKBTHTBYD33/18893J11",
		"/4iOhRr82034xTNd833kcfcpxpqee11723RMRIDKVBJHYQOAC33/18893d11",
		"/4Iotjw81001AlPk1363ebsqyzwhvs10521SZXAOEISQDZMGFQ2040632/18649v11",
		"/4JzahJ27343EEjZ935xetnxejlse2280BTMBVMDFURUCIDT171028/11452c12",
		"/4KhoMf82034azyt833ggdfqlkohe11723KDHGOEFEIJRSHHW33/18893a11",
		"/4lFutf82034semi833sskgobgdnh11723YLOTCOXDIRUULFC33/18893o11",
		"/4nquwV27343oWIr935bwodvfnbvs2280OGWQHNECJIGZYYS171028/11452n12",
		"/4OkWZA82021HuJD1469spvdnblqwy11723QTVYETNXNPEWBRV33/18893O11",
		"/4QqsTi27343kBQj935fumvyggkcz2280KBIHEKOQSRRNPRD171028/11452Q12",
		"/4QUQTy82034csAe833lwrsgiudlw11723KJLFLASMNJBUVBU33/18893l11",
		"/4sCtWa80791GpEG1195teuwrvsxbm8171MYTFABBBOQCNHXM13839/18625v22",
		"/4uIEqG80984Plej1363mrquqxbpan10521NFBEGXOICSJWGLT2040632/18649e11",
		"/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22",
		"/4WKKpm80770NJNV1619ryqbstiovr8171MWJCRLTEIYFSZON13839/18625R22",
		"/4YNHkz27343kuqq935gqlzyjuprt2280ZGSSUMJKWSWBHYZ171028/11452r12",
		"/4zpqUC82034DAGl833wksoseggvm11723RTEDZOQZANEJLUZ33/18893D11",
		"/4Ztukv82034WvDO833bszojdwlxh11723WOIYADVATHLZCHU33/18893T11",
		"/5BNYik82021ZCuc1469ohmzrdrimg11723ZEJPNQXXGADRZGP33/18893t11",
		"/5dWwxC80791dfpi1195dkmhnmcxcg8171ASVGXUTOSQPQBUC13839/18625c22",
		"/5eYbZo80702JGjN767ypgmbwniii8171AWRGFKLKLRPKLPM13839/18625X22",
		"/5fJPcu80770LrUx1619qorpaodvmu8171ALLOFCDQMYEKMIT13839/18625x22",
		"/5KVZii82034wxph833wvgqwyznux11723VKHWDAUHOVSIELI33/18893u11",
		"/5laNXF80790UzLG767fwqmqfwwsz8171JFPBMGWFWHOFAWE13839/18625r22",
		"/5NvKDQ80984riYX1363umvrmnphby10521FORIIULLATURFLG2040632/18649Z11",
		"/5TcdFH80792pQrq1362smjculxvav8171LYHBVZZVECVIUNG13839/18625M22",
		"/5ZsOzX27343NhJJ935juidljbzjg2280VCJIZLSVREFXRAQ171028/11452o12",
		"/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"/6DUEvB80790jvRD767claopcblok8171BPLORCVMSAKIGBA13839/18625i22",
		"/6jxIcg80770jiBE1619zosbmbjiag8171IJNXYRAQAPYTMBB13839/18625T22",
		"/6LhcWy82034HQMC833chpttjzuhj11723JHJFVULHXPTLFDZ33/18893t11",
		"/6RBIAx80791acSe1195cfqticzrpx8171RUIUVYDAHYVPUIZ13839/18625n22",
		"/6rUNVW82021CRkm1469qpolwswggz11723UNAUZEVABROYEXL33/18893k11",
		"/6UPeBZ80984EQhU1363qhpdnzbsew10521NBJBRWGNILCFBEZ2040632/18649i11",
		"/6yfdln80792fpEQ1362mzednzrstj8171YKZXKOVZXAHYVAE13839/18625N22",
		"/about",
		"/bqpuHJDME.js",
		"/brKcJyirt.js",
		"/buInbKFwT.js",
		"/favicon.ico",
		"/news",
		"/opt-out",
		"/privacy",
		"/robots.txt",
		"/t/robots.txt",
		"/track/3TMvbV80702vdqw767brbxkphbsk8171SFFDKUUQMREBQPG13839/18625p22",
		"/track/3ugKFR82021AvbF1469juczcbmecr11723CUTPQMVPKRFDQQG33/18893e11",
		"/wdb21",
		"/wdb36",
	}
	userAgents = []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36 (compatible; Google-InspectionTool/1.0;)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:128.0) Gecko/20100101 Firefox/128.0",
		"Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0",
	}
	templates = []string{
		"20195400242024-06-22 15:12:46{{.Padding}};326420845262644&isp=gmail.com&idesp=555",
		"20195400242024-06-22 15:12:46{{.Padding}};326420845262644&isp=icloud.com&idesp=555",
		"20195400242024-06-22 15:12:46{{.Padding}}@gmail.com;326420845262644&isp=gmail.com&idesp=555",
		"20195400242024-06-22 15:12:46{{.Padding}}@icloud.com;326420845262644&isp=icloud.com&idesp=555",
		"21512900242024-06-19 11:29:41{{.Padding}};327908039362531&isp=gmail.com&idesp=555",
		"21512900242024-06-19 11:29:41{{.Padding}};327908039362531&isp=icloud.com&idesp=555",
		"21512900242024-06-19 11:29:41{{.Padding}}@gmail.com;327908039362531&isp=gmail.com&idesp=555",
		"21512900242024-06-19 11:29:41{{.Padding}}@icloud.com;327908039362531&isp=icloud.com&idesp=555",
		"215131&{{.Padding}}@gmail.com&44337&ITA&2024-06-19 11:29:41&62531&3279080393&gmail.com",
		"215131&{{.Padding}}@icloud.com&44337&ITA&2024-06-19 11:29:41&62531&3279080393&icloud.com",
		"215131&{{.Padding}}&44337&ITA&2024-06-19 11:29:41&62531&3279080393&gmail.com",
		"215131&{{.Padding}}&44337&ITA&2024-06-19 11:29:41&62531&3279080393&icloud.com",
		"email={{.Padding}}@gmail.com&idNewsletter=44111&idioma=ITA",
		"email={{.Padding}}@icloud.com&idNewsletter=44111&idioma=ITA",
		"email={{.Padding}}&idNewsletter=44111&idioma=ITA",
	}
	httpWorker = make(chan spam, 1)
	stopped    = false
	counter    = uint64(0)
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for !stopped {
		select {

		case <-signalCh:
			stopped = true

		case s := <-httpWorker:
			go s.call()

		case httpWorker <- makeSpam():
			time.Sleep(time.Duration(counter/30) * time.Millisecond)
		}
	}
}
