package main

import (
	"errors"
	"flag"
	"log"
	"math"
	"regexp"
	"strings"
	"time"

	p "github.com/its-my-data/doubak/proto"
	"github.com/its-my-data/doubak/task"
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
var _ = flag.String(p.Flag_cookies_file.String(), "",
	"The input file path to the cookies.txt.")
var _ = flag.String(p.Flag_output_dir.String(), "./output",
	"The base output path.")
var incrementalRun = flag.Bool(p.Flag_incremental.String(), true,
	"Incremental or restart with override.")
var proxy = flag.String(p.Flag_proxy.String(), "",
	"Proxy to use when crawling.")
var numRetry = flag.Uint64(p.Flag_max_retry.String(), math.MaxUint64,
	"The number of retries when errors encountered.")
var defaultRequestDelay, _ = time.ParseDuration("100ms")
var requestDelay = flag.Duration(p.Flag_req_delay.String(), defaultRequestDelay,
	"Min time between any two requests, used to reduce server load. This may "+
		"be replaced by a QPS flag when proxy pool and parallel requests are implemented.")

func validateFlags() (tasks []string, categories []string, err error) {
	spaceRegex := regexp.MustCompile(`\s`)

	// Validate task list (order matters).
	strippedTasks := spaceRegex.ReplaceAllString(*tasksToRun, "")
	tasks = strings.Split(strippedTasks, ",")
	for _, t := range tasks {
		if _, ok := p.Task_value[t]; !ok {
			err = errors.New("unknown task name: " + t)
			return
		}
	}

	// Validate category list (order doesn't matter).
	strippedCategories := spaceRegex.ReplaceAllString(*targetCategories, "")
	categories = strings.Split(strippedCategories, ",")
	for _, c := range categories {
		if _, ok := p.Category_value[c]; !ok {
			err = errors.New("unknown category name: " + c)
			return
		}
	}

	return
}

func main() {
	flag.Parse()

	// Precheck flags that need preprocessing.
	log.Print("Validating flags... ")
	tasks, categories, parseErr := validateFlags()
	if parseErr != nil {
		log.Print("FAILED")
		log.Fatal(parseErr)
	} else {
		log.Print("PASS")
	}

	// Create selected tasks.
	taskMap := map[string]task.BaseInterface{}
	for _, t := range tasks {
		var taskImpl task.BaseInterface
		switch t {
		case p.Task_collect.String():
			taskImpl = task.NewCollector(userName, categories)
		case p.Task_parse.String():
			taskImpl = task.NewParser(categories)
		case p.Task_publish.String():
			taskImpl = task.NewPublisher(categories)
		}
		taskMap[t] = taskImpl
	}

	// Run the specific tasks' prechecks first.
	for taskName, t := range taskMap {
		log.Printf("Prechecking \"%s\"... ", taskName)
		if err := t.Precheck(); err != nil {
			log.Print("FAILED")
			log.Fatal(err)
		} else {
			log.Print("PASS")
		}
	}

	// Execute the tasks in input order.
	for _, taskName := range tasks {
		log.Printf("Running task \"%s\"... ", taskName)
		if err := taskMap[taskName].Execute(); err != nil {
			log.Printf("Task \"%s\" execution failed", taskName)
			log.Fatal(err)
		} else {
			log.Printf("Task \"%s\" competed", taskName)
		}
	}
}
