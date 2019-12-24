package parser

import (
	"bytes"
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

// RemoveExcludedLinks : Removes any links that contain any of the substrings to exclude
func RemoveExcludedLinks(links []string, exclude []string) []string {
	results := make([]string, 0)
	for _, link := range links {
		shouldInclude := true

		for _, subStr := range exclude {
			if strings.Contains(link, subStr) {
				shouldInclude = false
				break
			}
		}

		if shouldInclude {
			results = append(results, link)
		}
	}

	return results
}

// PrependDomainToLinks : Creates a new list of links with the given domain prepended
func PrependDomainToLinks(links []string, prefix string) []string {
	newLinks := make([]string, 0)
	for _, link := range links {
		var buffer bytes.Buffer
		buffer.WriteString(prefix)
		buffer.WriteString(link)
		newLinks = append(newLinks, buffer.String())
	}

	return newLinks
}

// TrimDocument : Trim all parts of the HTML document after a delimiter appears
func TrimDocument(htm string, delimiter string) string {
	parts := strings.Split(htm, delimiter)
	return parts[0]
}

func sanitizeLink(link string) string {
	newString := strings.ReplaceAll(link, "\"", "")
	return strings.ReplaceAll(newString, " ", "")
}
