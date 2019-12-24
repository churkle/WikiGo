package crawler

import (
	"WikiGo/parser"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Crawler : struct that has a source and destination page with a map cache of shortest distances
//          between the page (key) and destination
type Crawler struct {
	src        string
	dest       string
	domain     string
	pattern    []string
	exclude    []string
	trimMarker string
	limit      int
}

// NewCrawler : creates a new crawler object with src and dest pages
func NewCrawler(src string, dest string, domain string, pattern []string, exclude []string, trimMarker string, limit int) *Crawler {
	c := Crawler{src: src, dest: dest, domain: domain, pattern: pattern, exclude: exclude, trimMarker: trimMarker, limit: limit}
	return &c
}

// GetShortestPathToArticle : Takes two URLs and computes the shortest way to
//                           get from one to the other through links
func (c Crawler) GetShortestPathToArticle() ([]string, error) {
	for i := 1; i <= c.limit; i++ {
		history := make([]string, 0)
		result, err := c.crawl(c.src, history, i)

		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (c Crawler) crawl(url string, history []string, maxDepth int) (result []string, err error) {
	path := append(history, url)

	if url == c.dest {
		return path, nil
	}

	if len(path) > maxDepth {
		return nil, nil
	}

	if url == "" {
		return nil, errors.New("Empty URL")
	}

	htm, err := GetHTMLFromURL(url)
	if err != nil {
		return nil, err
	}

	htm = parser.TrimDocument(htm, c.trimMarker)

	links, err := parser.GetLinks(strings.NewReader(htm))
	links = parser.PrependDomainToLinks(parser.RemoveExcludedLinks(parser.FilterPatternLinks(links, c.pattern), c.exclude), c.domain)

	if err != nil {
		return nil, err
	}

	for _, link := range links {
		for _, site := range history {
			if link == site {
				return nil, nil
			}
		}

		result, err := c.crawl(link, path, maxDepth)

		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, nil
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
