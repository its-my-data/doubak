package task

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/its-my-data/doubak/proto"
	"github.com/its-my-data/doubak/util"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// TODO: use a separate library for URLs.

const DoubanURL = "https://www.douban.com/"
const PeopleURL = DoubanURL + "people/"

const startingPage = 1

var timePrefix = time.Now().Local().Format("20060102.1504")

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
	cu := util.NewColly()
	cu.OnError(func(r *colly.Response, err error) {
		exists = false
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
			//task.crawlBroadcastLists()
			task.crawlBroadcastDetail()
		case proto.Category_book.String():
			task.crawlBookLists()
		case proto.Category_movie.String():
			task.crawlMovieLists()
		case proto.Category_game.String():
			task.crawlGameLists()
		default:
			return errors.New("Category not implemented " + c)
		}
	}
	return nil
}

// crawlBroadcastLists downloads the list of broadcasts.
func (task *Collector) crawlBroadcastLists() error {
	page := startingPage
	q := util.NewQueue()
	c := util.NewColly()

	c.OnResponse(func(r *colly.Response) {
		fileName := fmt.Sprintf("%s_%s_p%d.html", timePrefix, proto.Category_broadcast, page)
		if err := task.saveResponse(r, fileName); err != nil {
			log.Println(err.Error())
		}

		body := string(r.Body)
		util.FailIfNeedLogin(&body)

		// Prepare for the next request.
		// Note that the number of broadcasts in each page somehow don't equal.
		// Therefore, I have to get at least an empty status page file.
		broadcastCount := strings.Count(body, "\"status-item\"")
		log.Println("Found", broadcastCount, "broadcasts/status.")
		if broadcastCount != 0 {
			page++
			url := PeopleURL + task.user + "/statuses?p=" + strconv.Itoa(page)
			q.AddURL(url)
			log.Printf("Added URL: %s. (Followed by sleeping.)\n", url)
			time.Sleep(util.RequestInterval)
		} else {
			log.Printf("All done with broadcast count %d (in page %d).\n", broadcastCount, page)
		}
	})

	// TODO: need a retry queue (either Requests, or go routines).
	q.AddURL(PeopleURL + task.user + "/statuses?p=" + strconv.Itoa(page))

	return q.Run(c)
}

// crawlBroadcastDetail downloads the detail of each broadcast by reading all downloaded broadcast lists.
func (task *Collector) crawlBroadcastDetail() error {
	fileNamePattern := fmt.Sprintf("*_%s_p*.html", proto.Category_broadcast)
	files := util.GetFilePathListWithPattern(task.outputDir, fileNamePattern)
	for _, fn := range files {
		log.Println("Found file:", fn)
		// TODO: finish this with goquery.
	}

	// TODO: handle each type of broadcasts.

	return errors.New("update the implementation")
}

func (task *Collector) crawlBookLists() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

func (task *Collector) crawlMovieLists() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

func (task *Collector) crawlGameLists() error {
	// TODO: update the implementation.
	return errors.New("update the implementation")
}

// TODO: implement more crawlers.

func (task *Collector) saveResponse(r *colly.Response, fileName string) error {
	fullPath := filepath.Join(task.outputDir, fileName)
	if err := r.Save(fullPath); err != nil {
		return err
	}
	log.Println("Saved", fullPath)
	return nil
}
