package task

import (
	"flag"
	"github.com/gocolly/colly/v2"
	p "github.com/its-my-data/doubak/proto"
	"log"
)

// Collector contains the information used by the collector.
type Collector struct {
	user       string
	categories []string
}

// NewCollector returns a new collector task and initialise it.
func NewCollector(categories []string) *Collector {
	return &Collector{
		user:       flag.Lookup(p.Flag_user.String()).Value.(flag.Getter).Get().(string),
		categories: categories,
	}
}

// Precheck validates the flags.
func (task *Collector) Precheck() error {
	// TODO: check user existance, etc.
	return nil
}

// Execute starts the collection.
func (task *Collector) Execute() error {
	// TODO: update the implementation.
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		log.Println("Found ULR: ", e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.Visit("http://douban.com/")
	return nil
}
