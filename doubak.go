package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"math"
	"time"
)

// Defining flags.
var userName = flag.String("user", "", "The Douban user name. e.g. mewcatcher")
var tasksToRun = flag.String("tasks", "collect, analyse, publish",
	"Tasks to run (order doesn't matter). Can be one/more of the following: "+
		"collect, analyse, publish.")
var targetCategories = flag.String("categories", "",
	"A comma separated content types list to crawl. Default is all. "+
		"Supported types are: book, movie, music, game, app, review.")
var outputDir = flag.String("output_dir", "./output", "The output path.")
var continueRun = flag.Bool("continue", true,
	"Continue or restart with override.")
var proxy = flag.String("proxy", "", "Proxy to use when crawling.")
var numRetry = flag.Uint64("max_retry", math.MaxUint64,
	"The number of retries when errors encountered.")
var defaultRequestDelay, _ = time.ParseDuration("100ms")
var requestDelay = flag.Duration("req_delay", defaultRequestDelay,
	"Delay betwee two requests, used to control QPS. This may be replaced by "+
		"a QPS flag when proxy pool and parallel requests are added.")

func main() {
	flag.Parse()

	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Found ULR: ", e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit("http://douban.com/")
}
