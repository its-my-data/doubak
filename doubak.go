package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/its-my-data/doubak/collector"
	p "github.com/its-my-data/doubak/proto"
	"math"
	"time"
)

// Defining flags.
var userName = flag.String(p.Flag_user.String(), "",
	"The Douban user name. e.g. mewcatcher")
var tasksToRun = flag.String(p.Flag_tasks.String(),
	p.ConcatProtoEnum(p.Task_name, ", "),
	"Tasks to run (order doesn't matter). Can be one/more of the following: "+
		p.ConcatProtoEnum(p.Task_name, ", ")+".")
var targetCategories = flag.String(p.Flag_categories.String(),
	p.ConcatProtoEnum(p.Category_name, ", "),
	"A comma separated content types list to crawl. Default is all. "+
		"Supported types are: "+p.ConcatProtoEnum(p.Category_name, ", ")+".")
var outputDir = flag.String(p.Flag_output_dir.String(), "./output",
	"The output path.")
var continueRun = flag.Bool(p.Flag_continue.String(), true,
	"Continue or restart with override.")
var proxy = flag.String(p.Flag_proxy.String(), "",
	"Proxy to use when crawling.")
var numRetry = flag.Uint64(p.Flag_max_retry.String(), math.MaxUint64,
	"The number of retries when errors encountered.")
var defaultRequestDelay, _ = time.ParseDuration("100ms")
var requestDelay = flag.Duration(p.Flag_req_delay.String(), defaultRequestDelay,
	"Min time between any two requests, used to reduce server load. This may "+
		"be replaced by a QPS flag when proxy pool and parallel requests are implemented.")

func main() {
	flag.Parse()

	collector.Collect()

	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Found ULR: ", e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit("http://douban.com/")
}
