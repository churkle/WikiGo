package crawler

import (
	"WikiGo/parser"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

//Crawler : struct that has a source and destination page with a map cache of shortest distances
//          between the page (key) and destination
type Crawler struct {
	src          string
	dest         string
	domain       string
	pattern      []string
	exclude      []string
	trimMarker   string
	limit        int
	shortestPath []string
	mux          sync.Mutex
}

// NewCrawler : creates a new crawler object with src and dest pages
func NewCrawler(src string, dest string, domain string, pattern []string, exclude []string, trimMarker string, limit int) *Crawler {
	c := Crawler{src: src, dest: dest, domain: domain, pattern: pattern, exclude: exclude, trimMarker: trimMarker, limit: limit}
	c.shortestPath = make([]string, 0)
	return &c
}

// GetShortestPathToArticle : Takes two URLs and computes the shortest way to
//                           get from one to the other through links
func (c *Crawler) GetShortestPathToArticle() ([]string, error) {
	for i := 1; i <= c.limit; i++ {
		history := make([]string, 0)
		var wg sync.WaitGroup

		wg.Add(1)
		c.crawl(c.src, history, i, &wg)

		wg.Wait()

		if c.shortestPath != nil && len(c.shortestPath) != 0 {
			fmt.Println("SUCCESS")
			printPath(c.shortestPath)
			return c.shortestPath, nil
		}
	}

	fmt.Println("FAIL")
	return nil, nil
}

func (c *Crawler) crawl(url string, history []string, maxDepth int, wg *sync.WaitGroup) {
	defer wg.Done()
	path := append(history, url)

	if url == c.dest {
		c.updateShortestPath(path)
		return
	}

	if len(path) > maxDepth {
		return
	}

	if url == "" {
		fmt.Println("Empty URL")
		return
	}

	htm, err := GetHTMLFromURL(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	htm = parser.TrimDocument(htm, c.trimMarker)

	links, err := parser.GetLinks(strings.NewReader(htm))
	links = parser.PrependDomainToLinks(parser.RemoveExcludedLinks(parser.FilterPatternLinks(links, c.pattern), c.exclude), c.domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range links {
		for _, site := range history {
			if link == site {
				return
			}
		}

		wg.Add(1)
		go c.crawl(link, path, maxDepth, wg)
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

//GetHTMLFromURL : Retrieves the HTML reader from a URL
func GetHTMLFromURL(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
