package util

import (
	"flag"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/its-my-data/doubak/proto"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// NewQueue creates the Colly Request Queue based on my needs.
func NewQueue() *queue.Queue {
	q, err := queue.New(
		1,                                           // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	if err != nil {
		log.Fatal(err)
	}
	return q
}

// NewColly creates the Colly Collector based on my needs.
func NewColly() *colly.Collector {
	// Calculate the cookies only during initialization.
	cookies := ""
	if f := flag.Lookup(proto.Flag_cookies_file.String()); f != nil && len(f.Value.String()) != 0 {
		c, err := LoadCookiesFileToString(f.Value.String())
		if err != nil {
			log.Println(err)
		}
		cookies = c
	}

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

		// Hacking the cookies.
		if len(cookies) != 0 {
			r.Headers.Set("Cookie", cookies)
		}

		r.Headers.Set("Referer", "https://www.douban.com/")
		r.Headers.Set("Host", "https://www.douban.com/")
	})

	c.SetRequestTimeout(5 * time.Minute)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Minute,
			KeepAlive: 5 * time.Minute,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       5 * time.Minute,
		TLSHandshakeTimeout:   5 * time.Minute,
		ExpectContinueTimeout: 5 * time.Minute,
	})
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		Delay:       5 * time.Second,
	})

	return c
}

// FailIfNeedLogin will fail the execution if the page requires login.
func FailIfNeedLogin(body *string) {
	if strings.Contains(*body, "<title>登录豆瓣</title>") {
		log.Fatal("Cannot continue, need to log in first")
	}
}
