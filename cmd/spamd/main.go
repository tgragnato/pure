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
	freshDomains = []string{
		"http://162.62.230.106",
		"http://kingdaros.com",
		"https://egsro.com",
		"https://solutionformen.com",
	}
	domains = []string{
		"http://162.62.230.106",
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
		"http://ashfordaikido.store",
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
		"http://kingdaros.com",
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
		"http://marypozas.com",
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
		"https://allcommunitiesonline.egsro.com",
		"https://bit.ly/3yM5GYw",
		"https://egsro.com",
		"https://gitlab.solutionformen.com",
		"https://jerprint.com",
		"https://m.egsro.com",
		"https://mail.egsro.com",
		"https://mail.solutionformen.com",
		"https://panelapi.mooo.com/tr.php",
		"https://solutionformen.com",
		"https://tinyurl.com/HIMALA12",
		"https://tinyurl.com/htsoni",
		"https://tinyurl.com/Loura41r",
		"https://tinyurl.com/Wolaf154",
		"https://www.egsro.com",
		"https://www.solutionformen.com",
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
		"/4BVbmS27343JzYU935yrxrokvdvf2280EWOQQAIFKKQKNGH171028/11452c12#0odw4f2i71nsop3zeomk",
		"/4HXjoo80791dYtM1195zuwviakgln8171EJFZXXVQGNPXAOK13839/18625R22#pdbppwwlaexditvthghaj",
		"/4IfNnV82034CQDT833zafgfloxof11723BMQLFBXKBTHTBYD33/18893J11#heeoxeiyblpllkghwvdcj",
		"/4iOhRr82034xTNd833kcfcpxpqee11723RMRIDKVBJHYQOAC33/18893d11#whczcljicohwgyrmuappj",
		"/4Iotjw81001AlPk1363ebsqyzwhvs10521SZXAOEISQDZMGFQ2040632/18649v11#qfhofadsqooaihokwkokk",
		"/4JzahJ27343EEjZ935xetnxejlse2280BTMBVMDFURUCIDT171028/11452c12#x5bkwog0lcexbya2jozj",
		"/4KhoMf82034azyt833ggdfqlkohe11723KDHGOEFEIJRSHHW33/18893a11#baudjimbzglzohchakxef",
		"/4lFutf82034semi833sskgobgdnh11723YLOTCOXDIRUULFC33/18893o11#rpstnsecwoifhdvlgpewv",
		"/4nquwV27343oWIr935bwodvfnbvs2280OGWQHNECJIGZYYS171028/11452n12#0qpqsh1e10bqzp7puz6h",
		"/4OkWZA82021HuJD1469spvdnblqwy11723QTVYETNXNPEWBRV33/18893O11#wcojolmmhlyxewtiadbbs",
		"/4QqsTi27343kBQj935fumvyggkcz2280KBIHEKOQSRRNPRD171028/11452Q12#njcpvi2sss799re8ad7u",
		"/4QUQTy82034csAe833lwrsgiudlw11723KJLFLASMNJBUVBU33/18893l11#drfagtjtlonnvyrfpheqf",
		"/4sCtWa80791GpEG1195teuwrvsxbm8171MYTFABBBOQCNHXM13839/18625v22#riwkuhwealeouqzihzvcb",
		"/4uIEqG80984Plej1363mrquqxbpan10521NFBEGXOICSJWGLT2040632/18649e11#uhrmkzosemksdrbqpzucd",
		"/4upvoz80770TqYK1619fcsahliicq8171PCLOYQTVNNOSXJT13839/18625S22#drsbyizemphcffpuxqwje",
		"/4WKKpm80770NJNV1619ryqbstiovr8171MWJCRLTEIYFSZON13839/18625R22#lhsglrclhqluewhayihdy",
		"/4YNHkz27343kuqq935gqlzyjuprt2280ZGSSUMJKWSWBHYZ171028/11452r12#4yn240fvh6e9sj6lrx05",
		"/4zpqUC82034DAGl833wksoseggvm11723RTEDZOQZANEJLUZ33/18893D11#jvqihqyiebzkihxapdrkw",
		"/4Ztukv82034WvDO833bszojdwlxh11723WOIYADVATHLZCHU33/18893T11#fomcumbeixzmtnknsawzh",
		"/5BNYik82021ZCuc1469ohmzrdrimg11723ZEJPNQXXGADRZGP33/18893t11#qbozugrwxkqrqyevizmav",
		"/5dWwxC80791dfpi1195dkmhnmcxcg8171ASVGXUTOSQPQBUC13839/18625c22#nfocppzyzkhclaigfdedr",
		"/5eYbZo80702JGjN767ypgmbwniii8171AWRGFKLKLRPKLPM13839/18625X22#gypdyeghnejaegotzgqyt",
		"/5fJPcu80770LrUx1619qorpaodvmu8171ALLOFCDQMYEKMIT13839/18625x22#xqxhaflofarvykfyvweyr",
		"/5KVZii82034wxph833wvgqwyznux11723VKHWDAUHOVSIELI33/18893u11#aymrkvzqzqbgsealcdbqi",
		"/5laNXF80790UzLG767fwqmqfwwsz8171JFPBMGWFWHOFAWE13839/18625r22#fbtzrkfnduqlqqiegsvcg",
		"/5NvKDQ80984riYX1363umvrmnphby10521FORIIULLATURFLG2040632/18649Z11#pfxglimjeauhwqekshuix",
		"/5TcdFH80792pQrq1362smjculxvav8171LYHBVZZVECVIUNG13839/18625M22#zuylgqjyjiflobmwjfcvk",
		"/5ZsOzX27343NhJJ935juidljbzjg2280VCJIZLSVREFXRAQ171028/11452o12",
		"/6Chwaz80702KKGq767hgugxsgkfh8171DAMTTNPHCPQMXUF13839/18625h22",
		"/6DUEvB80790jvRD767claopcblok8171BPLORCVMSAKIGBA13839/18625i22",
		"/6jxIcg80770jiBE1619zosbmbjiag8171IJNXYRAQAPYTMBB13839/18625T22",
		"/6LhcWy82034HQMC833chpttjzuhj11723JHJFVULHXPTLFDZ33/18893t11#phurrtaxpsezhmejnayaq",
		"/6RBIAx80791acSe1195cfqticzrpx8171RUIUVYDAHYVPUIZ13839/18625n22",
		"/6rUNVW82021CRkm1469qpolwswggz11723UNAUZEVABROYEXL33/18893k11#xssbmjceiisrtywcgffms",
		"/6UPeBZ80984EQhU1363qhpdnzbsew10521NBJBRWGNILCFBEZ2040632/18649i11",
		"/6yfdln80792fpEQ1362mzednzrstj8171YKZXKOVZXAHYVAE13839/18625N22#lqzxddadypshbybofumzp",
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
		"/track/3ugKFR82021AvbF1469juczcbmecr11723CUTPQMVPKRFDQQG33/18893e11#gkywchtcegrxwvnizpvuo",
		"/wdb21",
		"/wdb36",
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
		"Mozilla/5.0 (Windows NT 10.0; rv:129.0) Gecko/20100101 Firefox/129.0",
	}
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go iterate()
	for _, domain := range freshDomains {
		go func(domain string) {
			for {
				call(domain+generateRandomString(), userAgents[rand.Intn(len(userAgents))])
			}
		}(domain)
		go func(domain string) {
			for {
				call(domain+"/track"+generateRandomString(), userAgents[rand.Intn(len(userAgents))])
			}
		}(domain)
	}
	go func() {
		for range time.NewTicker(time.Minute).C {
			go iterate()
		}
	}()

	<-signalCh
}

func iterate() {
	items := []string{}
	randomItems := []string{}
	randomTracks := []string{}
	for _, domain := range domains {
		for _, path := range paths {
			items = append(items, domain+path)
			randomItems = append(randomItems, domain+generateRandomString())
			randomTracks = append(randomTracks, domain+"/track"+generateRandomString())
		}
	}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	for index, url := range items {
		go call(url, userAgents[index%len(userAgents)])
		go call(randomItems[index], userAgents[rand.Intn(len(userAgents))])
		go call(randomTracks[index], userAgents[rand.Intn(len(userAgents))])
		time.Sleep(time.Minute / time.Duration(len(items)))
	}
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

func generateRandomString() string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := rand.Intn(10) + 60
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return fmt.Sprintf("/%s/%s#%s", result[:length-30], result[length-30:length-22], result[length-22:])
}
