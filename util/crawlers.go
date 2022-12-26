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

const RequestTimeout = 5 * time.Minute
const RequestInterval = 3 * time.Second

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

	c.OnError(func(r *colly.Response, err error) {
		t := string(r.Body)
		if strings.Contains(t, "页面不存在") {
			// Real user not found.
		} else {
			// Usually caused by timeouts, but just in case.
			log.Println("Request URL:", r.Request.URL, "\nError:", err)
			log.Println("Unknown error with request body: ", string(r.Body))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

		// Hacking the cookies.
		if len(cookies) != 0 {
			r.Headers.Set("Cookie", cookies)
		}

		r.Headers.Set("Referer", "https://www.douban.com/")
		r.Headers.Set("Host", "https://www.douban.com/")
	})

	c.SetRequestTimeout(RequestTimeout)
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   RequestTimeout,
			KeepAlive: RequestTimeout,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       RequestTimeout,
		TLSHandshakeTimeout:   RequestTimeout,
		ExpectContinueTimeout: RequestTimeout,
	})
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		Delay:       RequestInterval, // FIXME: this doesn't seem to do anything?
	})

	return c
}

// FailIfNeedLogin will fail the execution if the page requires login.
func FailIfNeedLogin(body *string) {
	if strings.Contains(*body, "<title>登录豆瓣</title>") {
		log.Fatal("Cannot continue, need to log in first")
	}
}
