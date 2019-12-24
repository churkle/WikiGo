package crawler

import (
	"WikiGo/parser"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

//Crawler : struct that has a source and destination page with a map cache of shortest distances
//          between the page (key) and destination
type Crawler struct {
	src   string
	dest  string
	limit int
}

// NewCrawler : creates a new crawler object with src and dest pages
func NewCrawler(src string, dest string, limit int) *Crawler {
	c := Crawler{src: src, dest: dest, limit: limit}
	return &c
}

// GetShortestPathToArticle : Takes two URLs and computes the shortest way to
//                           get from one to the other through links
func (c Crawler) GetShortestPathToArticle() ([]string, error) {
	for i := 1; i <= c.limit; i++ {
		result, err := c.crawl(c.src, make([]string, 0), c.limit)

		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

func (c Crawler) crawl(url string, history []string, maxDepth int) (path []string, err error) {
	if len(history)-1 > maxDepth {
		return nil, nil
	}

	if url == "" {
		return nil, errors.New("Empty URL")
	}

	htmlReader, err := GetHTMLReaderFromURL(url)
	if err != nil {
		return nil, err
	}

	links, err := parser.GetLinks(strings.NewReader(htmlReader))
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		path := append(history, url)
		if link == c.dest {
			return path, nil
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

//GetHTMLReaderFromURL : Retrieves the HTML reader from a URL
func GetHTMLReaderFromURL(url string) (string, error) {
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
