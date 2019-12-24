package parser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

//GetLinks : Parse all links from the HTML document
func GetLinks(htm io.Reader) ([]string, error) {
	htmlTree, err := html.Parse(htm)
	links := make(map[string]bool)

	if err != nil {
		return nil, err
	}

	var crawl func(node *html.Node)
	crawl = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			attributes := node.Attr
			for _, attr := range attributes {
				if attr.Key == "href" {
					links[attr.Val] = true
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawl(child)
		}
	}

	crawl(htmlTree)

	keys := make([]string, 0, len(links))
	for key := range links {
		keys = append(keys, key)
	}
	return keys, nil
}

//FilterPatternLinks : This will filter out links without the prefix
func FilterPatternLinks(links []string, patterns []string) []string {
	resultLinks := make([]string, 0)

	for _, link := range links {
		for _, pattern := range patterns {
			if strings.HasPrefix(link, pattern) {
				resultLinks = append(resultLinks, link)
				break
			}
		}
	}

	return resultLinks
}

func sanitizeLink(link string) string {
	newString := strings.ReplaceAll(link, "\"", "")
	return strings.ReplaceAll(newString, " ", "")
}
