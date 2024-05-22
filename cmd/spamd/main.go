package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
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
	}
	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
	}
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < numberOfWorkers; i++ {
		go func(wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc, i int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-signalCh:
					cancel()
					return
				default:
					for _, url := range urls {
						req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						req.Header.Set("User-Agent", userAgents[i%len(userAgents)])
						resp, err := httpClient.Do(req)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						resp.Body.Close()
					}
				}
			}
		}(&wg, ctx, cancel, i)
	}
	wg.Wait()
}
