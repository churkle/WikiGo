package parser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

//GetLinks : Parse all links from the HTML document
func GetLinks(htm io.Reader) ([]string, error) {
	htmlTree, err := html.Parse(htm)
	links := make([]string, 0)

	if err != nil {
		return nil, err
	}

	var crawl func(node *html.Node)
	crawl = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			attributes := node.Attr
			for _, attr := range attributes {
				if attr.Key == "href" {
					links = append(links, sanitizeLink(attr.Val))
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawl(child)
		}
	}

	crawl(htmlTree)
	return links, nil
}

func sanitizeLink(link string) string {
	newString := strings.ReplaceAll(link, "\"", "")
	return strings.ReplaceAll(newString, " ", "")
}
