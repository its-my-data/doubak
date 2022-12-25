package task

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/its-my-data/doubak/proto"
	"github.com/its-my-data/doubak/util"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TODO: use a separate library for URLs.

const DoubanURL = "https://www.douban.com/"
const PeopleURL = DoubanURL + "people/"

// Collector contains the information used by the collector.
type Collector struct {
	user       string
	categories []string
	outputDir  string
}

// NewCollector returns a new collector task and initialise it.
func NewCollector(userName *string, categories []string) *Collector {
	return &Collector{
		user:       *userName,
		categories: categories,
	}
}

// Precheck validates the flags.
func (task *Collector) Precheck() error {
	// Initialize the top most directory for Collector.
	if path, err := util.GetPathWithCreation(util.CollectorPathPrefix); err != nil {
		return err
	} else {
		task.outputDir = path
	}
	log.Println("New output path saved:", task.outputDir)

	// Check user existence.
	exists := true
	cu := colly.NewCollector()
	// TODO: remove.
	//cu.OnResponse(func(r *colly.Response) {
	//	log.Println(string(r.Body))
	//})
	//cu.OnHTML("li", func(e *colly.HTMLElement) {
	//	log.Println(e.Text)
	//})
	cu.SetRequestTimeout(5 * time.Minute)
	cu.OnError(func(r *colly.Response, err error) {
		exists = false
		t := string(r.Body)
		if strings.Contains(t, "页面不存在") {
			// Real user not found.
		} else {
			// Usually caused by timeouts, but just in case.
			log.Println("Request URL:", r.Request.URL, "\nError:", err)
			log.Println("Unknown error with request body: ", string(r.Body))
		}
	})

	// Error handled separately.
	_ = cu.Visit(PeopleURL + task.user + "/")

	if !exists {
		return errors.New("user '" + task.user + "' does not exist")
	}
	return nil
}

// Execute starts the collection.
func (task *Collector) Execute() error {
	for _, c := range task.categories {
		switch c {
		case proto.Category_broadcast.String():
			task.crawlBroadcasts()
		case proto.Category_book.String():
			task.crawlBooks()
		case proto.Category_movie.String():
			task.crawlMovies()
		case proto.Category_game.String():
			task.crawlGames()
		default:
			return errors.New("Category not implemented " + c)
		}
	}
	return nil
}

func (task *Collector) crawlBroadcasts() error {
	// TODO: need to be global.
	timePrefix := time.Now().Local().Format("20060102.1504")
	page := 1

	q, _ := queue.New(
		1,                                           // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c := colly.NewCollector(
		colly.MaxDepth(1),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response header:", r.Headers)

		fileName := fmt.Sprintf("%s_broadcast_p%d.html", timePrefix, page)
		fullPath := filepath.Join(task.outputDir, fileName)

		if err := r.Save(fullPath); err != nil {
			log.Println(err.Error())
		}
		log.Println("Saved", fullPath)

		body := string(r.Body)
		if strings.Contains(body, "<title>登录豆瓣</title>") {
			log.Fatal("Cannot continue, need to log in first")
		}

		// Prepare for the next request.
		broadcastCount := strings.Count(body, "\"status-item\"")
		if broadcastCount == 20 {
			page++
			url := PeopleURL + task.user + "/statuses?p=" + strconv.Itoa(page)
			q.AddURL(url)
			log.Printf("Added URL: %s. (Followed by sleeping.)\n", url)
			time.Sleep(3 * time.Second)
		} else {
			log.Printf("All done with broadcast count %d (in page %d).\n", broadcastCount, page)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)

		// Hacking the cookies.
		cookiesFlag := flag.Lookup(proto.Flag_cookies_file.String())
		if cookiesFlag != nil && len(cookiesFlag.Value.String()) != 0 {
			cookies, err := util.LoadCookiesFileToString(cookiesFlag.Value.String())
			if err != nil {
				log.Println(err)
			}
			r.Headers.Set("Cookie", cookies)
		}

		r.Headers.Set("Referer", "https://www.douban.com/")
		r.Headers.Set("Host", "https://www.douban.com/")
	})

	// TODO: move the creator to a separate file so that we can set session and user-agent.
	// TODO: set cookies.
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

	// TODO: need a retry queue (either Requests, or go routines).
	q.AddURL(PeopleURL + task.user + "/statuses?p=" + strconv.Itoa(page))

	return q.Run(c)
}

func (task *Collector) crawlBooks() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

func (task *Collector) crawlMovies() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

func (task *Collector) crawlGames() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

// TODO: implement more crawlers.
