package task

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"log"
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
	// TODO: update the implementation.
	c := colly.NewCollector()
	// TODO: remove.
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	log.Println("Found ULR: ", e.Attr("href"))
	//})
	//c.OnRequest(func(r *colly.Request) {
	//	log.Println("Visiting", r.URL)
	//})
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		Delay:       5 * time.Second,
	})
	// TODO: need a retry queue (either Requests, or go routines).
	c.Visit("https://douban.com/")
	return nil
}
