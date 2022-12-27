package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/its-my-data/doubak/proto"
	"github.com/its-my-data/doubak/util"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TODO: use a separate library for URLs.

const DoubanURL = "https://www.douban.com/"
const MovieURL = "https://movie.douban.com/"
const BookURL = "https://book.douban.com/"
const MusicURL = "https://music.douban.com/"

const PeopleURL = DoubanURL + "people/"
const MoviePeopleURL = MovieURL + "people/"
const BookPeopleURL = BookURL + "people/"
const MusicPeopleURL = MusicURL + "people/"

const startingPage = 1
const startingItemId = 0

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
		log.Println("New output path saved:", task.outputDir)

		// Initialize other directories.
		p, _ := util.GetPathWithCreation(util.CollectorPathPrefix + util.ItemPathPrefix)
		log.Println("Created path:", p)
	}

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
// Broadcast should be synced regularly as it tracks all other category changes;
// although, it can still miss something that is not broadcast.
// In this case, we still want to run the sync on other categories on a less frequent basis.
// TODO: handle errors.
func (task *Collector) Execute() error {
	for _, c := range task.categories {
		switch c {
		case proto.Category_broadcast.String():
			task.crawlBroadcastLists()
			task.crawlBroadcastDetail()
		case proto.Category_book.String():
			task.crawlBookListDispatcher()
			task.crawlItemDetails(proto.Category_book, "li.subject-item > div.info > h2 > a")
		case proto.Category_movie.String():
			task.crawlMovieListDispatcher()
			// TODO: collect each movie details.
		case proto.Category_game.String():
			task.crawlGameListDispatcher()
			task.crawlItemDetails(proto.Category_game, "div.common-item > div.content > div.title > a:nth-child(1)")
		case proto.Category_music.String():
			task.crawlMusicListDispatcher()
			task.crawlItemDetails(proto.Category_music, "div.item > div.info > ul > li.title > a:nth-child(1)")
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
		// "p" means page.
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
	// Known data types are exhaust types of broadcast items.
	knownTypes := map[string]int{
		"game":    0, // A game.
		"movie":   0, // A movie.
		"book":    0, // A book.
		"music":   0, // A music album.
		"sns":     0, // A micro-post (might with pictures) that "someone made". If it's not made by me, it will have a different user ID in URL like (https://www.douban.com/people/54763828/status/3805173111/).
		"app":     0, // (Unsupported) An app.
		"ilmen":   0, // (Unsupported) An item that "I liked" (essentially a picture and a title).
		"fav":     0, // (Unsupported) A post/diary that "I liked".
		"olivia":  0, // (Unsupported) A discussion thread that "I participated" (sample: https://movie.douban.com/subject/3243582/discussion/637265942/).
		"doulist": 0, // (Unsupported) A douban list (sample https://douc.cc/0NnLLT or https://www.douban.com/game/25892303/).
		"rec":     0, // (Unsupported) A discussion thread that "I recommended" (https://douc.cc/4uYky6 or https://www.douban.com/group/topic/12410327/?_i=2003765TN2GQHs).
		"board":   0, // (Unsupported) A deprecated type that describes participating a movie pooling event.
		"":        0, // (Unsupported) A re-share of something.
	}

	fileNamePattern := fmt.Sprintf("*_%s_p*.html", proto.Category_broadcast)
	files := util.GetFilePathListWithPattern(task.outputDir, fileNamePattern)
	for _, fn := range files {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(util.ReadEntireFile(fn)))
		if err != nil {
			log.Println("Error reading", fn, "with message", err)
		}

		doc.Find("div.status-wrapper > div.status-item").Each(func(_ int, sel *goquery.Selection) {
			dataType := sel.AttrOr("data-target-type", "unspecified")
			if _, ok := knownTypes[dataType]; ok {
				// Do some statistics.
				knownTypes[dataType]++
			} else {
				html, _ := sel.Html()
				log.Printf("[WARNING] Found broadcast of type \"%s\" in %s\nFull element:\n%s\n", dataType, fn, strings.TrimSpace(html))
				// For statistical purpose and avoiding log spamming.
				knownTypes[dataType] = 1
			}
		})
	}

	// Pretty print the statistics.
	if b, err := json.MarshalIndent(knownTypes, "", "  "); err != nil {
		log.Fatal("error:", err)
	} else {
		log.Println("Statistics:", string(b))
	}

	// TODO: handle each type of broadcasts.

	return nil
}

func (task *Collector) crawlBookListDispatcher() error {
	// The book entry (https://book.douban.com/people/<user_name>/) contains the following parts:
	// - Read books.
	// - To-read books.
	// - Reading books.
	// - Others. (Not supported.)

	// Book list starts with item ID (which is 0). Each page has 15 items.
	// https://book.douban.com/people/mewcatcher/collect?start=<ID>&sort=time&rating=all&filter=all&mode=grid
	nRead := 0
	nToRead := 0
	nReading := 0
	c := util.NewColly()
	c.OnHTML("div#db-book-mine > div > h2", func(e *colly.HTMLElement) {
		secText := e.Text
		re := regexp.MustCompile("[0-9]+")
		nParsed, _ := strconv.Atoi(re.FindString(secText))

		switch {
		case strings.Contains(secText, "在读"):
			nReading = nParsed
			log.Println("Found reading books:", nReading)
		case strings.Contains(secText, "读过"):
			nRead = nParsed
			log.Println("Found read books:", nRead)
		case strings.Contains(secText, "想读"):
			nToRead = nParsed
			log.Println("Found to-read books:", nToRead)
		default:
			log.Println("Ignoring:", util.MergeSpaces(&secText))
		}
	})
	c.Visit(BookPeopleURL + task.user + "/")

	if err := task.crawlBookLists(nRead, "read", "collect"); err != nil {
		return err
	}
	if err := task.crawlBookLists(nToRead, "toread", "wish"); err != nil {
		return err
	}
	if err := task.crawlBookLists(nReading, "reading", "do"); err != nil {
		return err
	}
	return nil
}

func (task *Collector) crawlBookLists(totalItems int, tag string, urlAction string) error {
	const pageStep = 15
	urlTemplate := fmt.Sprintf("https://book.douban.com/people/%s/%s?start=%%d&sort=time&rating=all&filter=all&mode=grid", task.user, urlAction)
	return task.crawlItemLists(proto.Category_book, totalItems, pageStep, tag, urlTemplate)
}

func (task *Collector) crawlMovieListDispatcher() error {
	// The movie entry (https://movie.douban.com/people/<user_name>/) contains the following parts:
	// - Watched movies.
	// - To-watch movies.
	// - Watching movies.
	// - Favorite actors. (Not supported.)
	// - Movie Q&A. (Not supported.)

	// Movie list starts with item ID (which is 0). Each page has 15 items.
	// https://movie.douban.com/people/mewcatcher/collect?start=<ID>&sort=time&rating=all&filter=all&mode=grid
	nWatched := 0
	nToWatch := 0
	nWatching := 0
	c := util.NewColly()
	c.OnHTML("div#db-movie-mine > h2", func(e *colly.HTMLElement) {
		secText := e.Text
		re := regexp.MustCompile("[0-9]+")
		nParsed, _ := strconv.Atoi(re.FindString(secText))

		switch {
		case strings.Contains(secText, "看过"):
			nWatched = nParsed
			log.Println("Found watched movies:", nWatched)
		case strings.Contains(secText, "想看"):
			nToWatch = nParsed
			log.Println("Found to-watch movies:", nToWatch)
		case strings.Contains(secText, "在看"):
			nWatching = nParsed
			log.Println("Found watching movies:", nWatching)
		default:
			log.Println("Ignoring:", util.MergeSpaces(&secText))
		}
	})
	c.Visit(MoviePeopleURL + task.user + "/")

	if err := task.crawlMovieLists(nWatched, "watched", "collect"); err != nil {
		return err
	}
	if err := task.crawlMovieLists(nToWatch, "towatch", "wish"); err != nil {
		return err
	}
	if err := task.crawlMovieLists(nWatching, "watching", "do"); err != nil {
		return err
	}
	return nil
}

func (task *Collector) crawlMovieLists(totalItems int, tag string, urlAction string) error {
	const pageStep = 15
	urlTemplate := fmt.Sprintf("https://movie.douban.com/people/%s/%s?start=%%d&sort=time&rating=all&filter=all&mode=grid", task.user, urlAction)
	return task.crawlItemLists(proto.Category_movie, totalItems, pageStep, tag, urlTemplate)
}

func (task *Collector) crawlGameListDispatcher() error {
	// The game page does not have an entry (https://www.douban.com/people/<user_name>/games?action=<action>).
	// However, each page contains the following parts:
	// - To-play games.
	// - Playing games.
	// - Played games.
	// - Liked games. (Not supported due to legacy issue.)

	// Game list starts with item ID (which is 0). Each page has 15 items.
	// https://www.douban.com/people/mewcatcher/games?action=wish&start=<ID>
	nToPlay := 0
	nPlaying := 0
	nPlayed := 0
	c := util.NewColly()
	c.OnHTML("div.article > div.tabs > a", func(e *colly.HTMLElement) {
		secText := e.Text
		re := regexp.MustCompile("[0-9]+")
		nParsed, _ := strconv.Atoi(re.FindString(secText))

		switch {
		case strings.Contains(secText, "想玩"):
			nToPlay = nParsed
			log.Println("Found to-play games:", nToPlay)
		case strings.Contains(secText, "在玩"):
			nPlaying = nParsed
			log.Println("Found playing games:", nPlaying)
		case strings.Contains(secText, "玩过"):
			nPlayed = nParsed
			log.Println("Found played games:", nPlayed)
		default:
			log.Println("Ignoring:", util.MergeSpaces(&secText))
		}
	})
	c.Visit(PeopleURL + task.user + "/games")

	if err := task.crawlGameLists(nPlayed, "played", "collect"); err != nil {
		return err
	}
	if err := task.crawlGameLists(nToPlay, "toplay", "wish"); err != nil {
		return err
	}
	if err := task.crawlGameLists(nPlaying, "playing", "do"); err != nil {
		return err
	}
	return nil
}

func (task *Collector) crawlGameLists(totalItems int, tag string, urlAction string) error {
	const pageStep = 15
	urlTemplate := fmt.Sprintf("https://www.douban.com/people/%s/games?action=%s&start=%%d", task.user, urlAction)
	return task.crawlItemLists(proto.Category_game, totalItems, pageStep, tag, urlTemplate)
}

func (task *Collector) crawlMusicListDispatcher() error {
	// The music entry (https://music.douban.com/people/<user_name>/) contains the following parts:
	// - Listened album.
	// - To-listen album.
	// - Listening album.
	// - Others. (Not supported.)

	// Album music list starts with item ID (which is 0). Each page has 15 items.
	// https://music.douban.com/people/mewcatcher/collect?start=<ID>&sort=time&rating=all&filter=all&mode=grid
	nListened := 0
	nToListen := 0
	nListening := 0
	c := util.NewColly()
	c.OnHTML("div#db-music-mine > h2", func(e *colly.HTMLElement) {
		secText := e.Text
		re := regexp.MustCompile("[0-9]+")
		nParsed, _ := strconv.Atoi(re.FindString(secText))

		switch {
		case strings.Contains(secText, "在听"):
			nListening = nParsed
			log.Println("Found listening album music:", nListening)
		case strings.Contains(secText, "听过"):
			nListened = nParsed
			log.Println("Found listened album music:", nListened)
		case strings.Contains(secText, "想听"):
			nToListen = nParsed
			log.Println("Found to-listen album music:", nToListen)
		default:
			log.Println("Ignoring:", util.MergeSpaces(&secText))
		}
	})
	c.Visit(MusicPeopleURL + task.user + "/")

	if err := task.crawlMusicLists(nListened, "listened", "collect"); err != nil {
		return err
	}
	if err := task.crawlMusicLists(nToListen, "tolisten", "wish"); err != nil {
		return err
	}
	if err := task.crawlMusicLists(nListening, "listening", "do"); err != nil {
		return err
	}
	return nil
}

func (task *Collector) crawlMusicLists(totalItems int, tag string, urlAction string) error {
	const pageStep = 15
	urlTemplate := fmt.Sprintf("https://music.douban.com/people/%s/%s?start=%%d&sort=time&rating=all&filter=all&mode=grid", task.user, urlAction)
	return task.crawlItemLists(proto.Category_music, totalItems, pageStep, tag, urlTemplate)
}

// TODO: implement more crawlers.

// crawlItemLists downloads an item list universally.
func (task *Collector) crawlItemLists(cat proto.Category, totalItems int, pageStep int, tag string, urlTemplate string) error {
	// Validations.
	if strings.Count(urlTemplate, "%d") != 1 {
		return errors.New("URL template should have exact one %d placeholder")
	}

	startingItem := startingItemId
	c := util.NewColly()

	c.OnResponse(func(r *colly.Response) {
		// "l" means list.
		fileName := fmt.Sprintf("%s_%s_%s_l%d-%d.html", timePrefix, cat, tag, startingItem, startingItem+pageStep)
		if err := task.saveResponse(r, fileName); err != nil {
			log.Println(err.Error())
		}

		body := string(r.Body)
		util.FailIfNeedLogin(&body)

		itemCount := strings.Count(body, task.getItemMatcherPattern(cat))
		log.Println("Found", itemCount, cat.String()+"(s).")
		if itemCount != pageStep {
			log.Printf("Potential last %s page reached with count %d (in file %s).\n", cat, itemCount, fileName)
		}
	})

	for ; startingItem < totalItems; startingItem += pageStep {
		// TODO: implement retry strategy and incremental strategy.
		time.Sleep(util.RequestInterval)
		// Note that URL template should have a "%d" placeholder.
		url := fmt.Sprintf(urlTemplate, startingItem)
		err := c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (task *Collector) crawlItemDetails(cat proto.Category, selector string) error {
	q := util.NewQueue()
	inputFileNamePattern := fmt.Sprintf("*_%s_*.html", cat)
	files := util.GetFilePathListWithPattern(task.outputDir, inputFileNamePattern)
	for _, fn := range files {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(util.ReadEntireFile(fn)))
		if err != nil {
			log.Println("Error reading", fn, "with message", err)
		}

		// TODO: handle subject-item missing div.info issue.
		doc.Find(selector).Each(func(_ int, sel *goquery.Selection) {
			url, exists := sel.Attr("href")
			if !exists {
				log.Fatal("Found item without link", sel.Text())
			}

			// TODO: handle incremental option to check local file exists.
			q.AddURL(url)
		})
	}

	size, _ := q.Size()
	log.Println("Detail queue size is:", size)

	count := 1
	c := util.NewColly()
	c.OnResponse(func(r *colly.Response) {
		// Hack around to continue progress.
		if count < 0 {
			return
		}

		// Extract ID from response (using the first occurrence of number string).
		re := regexp.MustCompile("[0-9]+")
		id, _ := strconv.Atoi(re.FindString(r.Request.URL.String()))

		fileName := fmt.Sprintf("%s_%s_%d.html", timePrefix, cat, id)
		if err := task.saveResponse(r, util.ItemPathPrefix+fileName); err != nil {
			log.Println(err.Error())
		}

		body := string(r.Body)
		util.FailIfNeedLogin(&body)

		log.Println("Progress", count, "/", size)
		count++

		// TODO: replace this with a proper rate limiter.
		time.Sleep(util.RequestInterval)
	})
	return q.Run(c)
}

// getItemMatcherPattern returns the string matcher to identify matched item in the page.
func (task *Collector) getItemMatcherPattern(cat proto.Category) string {
	switch cat {
	case proto.Category_book:
		return "class=\"subject-item\""
	case proto.Category_movie, proto.Category_music:
		return "class=\"item\""
	case proto.Category_game:
		return "class=\"common-item\""
	default:
		log.Fatal("Nothing to match for category", cat)
		return "!!nothing-to-match!!"
	}
}

func (task *Collector) saveResponse(r *colly.Response, fileName string) error {
	fullPath := filepath.Join(task.outputDir, fileName)
	if err := r.Save(fullPath); err != nil {
		return err
	}
	log.Println("Saved", fullPath)
	return nil
}
