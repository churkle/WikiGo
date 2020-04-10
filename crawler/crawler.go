package crawler

import (
	"WikiGo/db"
	"WikiGo/parser"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// FileLimit : max number of files that can be open at once
const (
	FileLimit = 1000
)

// Crawler : struct that has a source and destination page with a map cache of shortest distances
//          between the page (key) and destination
type Crawler struct {
	wikiParser    *parser.Parser
	src           string
	dest          string
	srcTitle      string
	destTitle     string
	limit         int
	isWebCrawler  bool
	shortestPath  []string
	netClient     http.Client
	mux           sync.Mutex
	dbService     *db.Service
	adjacencyList map[string][]string
	urlMap        map[string]string
}

// NewCrawler : creates a new Crawler object with src and dest pages
func NewCrawler(src string, dest string, domain string, pattern []string, exclude []string, trimMarker []string,
	limit int, isWebCrawler bool, dbService *db.Service) *Crawler {

	c := Crawler{src: src, dest: dest, limit: limit, isWebCrawler: isWebCrawler, dbService: dbService}
	if c.dbService != nil {
		c.adjacencyList = c.dbService.GetPageGraph()
	}
	c.wikiParser = parser.NewParser(domain, pattern, exclude, trimMarker)
	c.shortestPath = make([]string, 0)
	srcHTML, _ := c.getHTMLFromURL(src)
	destHTML, _ := c.getHTMLFromURL(dest)
	c.srcTitle, _ = c.wikiParser.ExtractDocumentTitle(srcHTML)
	c.destTitle, _ = c.wikiParser.ExtractDocumentTitle(destHTML)
	tr := &http.Transport{
		MaxIdleConns:        15,
		MaxIdleConnsPerHost: 15,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true,
	}
	c.netClient = http.Client{Transport: tr}
	return &c
}

// GetShortestPathToArticle : Takes two URLs and computes the shortest way to
//                           get from one to the other through links
func (c *Crawler) GetShortestPathToArticle() ([]string, error) {
	if c.srcTitle == "" || c.destTitle == "" {
		return nil, errors.New("Unable to retrieve src or destination page")
	}

	for i := 1; i <= c.limit; i++ {
		if c.shortestPath != nil && len(c.shortestPath) != 0 {
			fmt.Println("SUCCESS")
			printPath(c.shortestPath)
			return c.shortestPath, nil
		}

		fmt.Println(i)
		history := make([]string, 0)
		maxChan := make(chan bool, FileLimit)
		var wg sync.WaitGroup

		wg.Add(1)
		maxChan <- true
		c.crawl(c.src, history, i, &wg, maxChan)

		wg.Wait()
	}

	if c.shortestPath != nil && len(c.shortestPath) != 0 {
		fmt.Println("SUCCESS")
		printPath(c.shortestPath)
		return c.shortestPath, nil
	}

	fmt.Println("FAIL")
	return nil, nil
}

func (c *Crawler) crawl(url string, history []string, maxDepth int, wg *sync.WaitGroup, maxChan chan bool) {
	defer wg.Done()
	defer func(maxChan chan bool) { <-maxChan }(maxChan)

	var title string
	var links []string

	if c.urlMap != nil {
		title = c.urlMap[url]
	}

	if title == "" {
		htm, err := c.getHTMLFromURL(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		title, err = c.wikiParser.ExtractDocumentTitle(htm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := append(history, title)
	fmt.Println(title)

	if title == c.destTitle {
		c.updateShortestPath(path)
		return
	}

	if len(path) > maxDepth {
		return
	}

	if title == "" {
		fmt.Println("Error retrieving document")
		return
	}

	if c.adjacencyList != nil {
		links = c.adjacencyList[title]
	}

	if links == nil {
		htm, err := c.getHTMLFromURL(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		links, _ = c.wikiParser.GetLinks(htm)
	}

	for _, link := range links {
		linkTitle, err := c.wikiParser.ExtractDocumentTitle(link)
		if err != nil {
			fmt.Println(err)
			continue
		}
		visited := false

		for _, site := range history {
			if linkTitle == site {
				visited = true
			}
		}

		if !visited {
			wg.Add(1)
			maxChan <- true
			go c.crawl(link, path, maxDepth, wg, maxChan)
		}
	}

	return
}

func (c *Crawler) updateShortestPath(path []string) {
	c.mux.Lock()
	if c.shortestPath != nil && len(c.shortestPath) != 0 {
		if len(path) < len(c.shortestPath) {
			c.shortestPath = path
		}
	} else {
		c.shortestPath = path
	}
	c.mux.Unlock()
}

func printPath(path []string) {
	for _, site := range path {
		fmt.Print(site + " -> ")
	}
	fmt.Println("")
}

func (c *Crawler) getHTMLFromURL(url string) (string, error) {
	var result []byte
	var err error
	if c.isWebCrawler {
		resp, err := c.netClient.Get(url)

		if err != nil {
			return "", err
		}

		defer resp.Body.Close()
		result, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}
	} else {
		result, err = ioutil.ReadFile(url)

		if err != nil {
			return "", err
		}
	}

	return string(result), nil
}
