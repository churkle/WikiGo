package parser

import (
	"bytes"
	"golang.org/x/net/html"
	"strings"
)

// Parser : struct that takes parses HTML documents, using a domain to find links for,
//          patterns to look for, links to exclude, and a trim marker to crop HTML at
type Parser struct {
	domain      string
	pattern     []string
	exclude     []string
	trimMarkers []string
}

// NewParser : Creates a new parser object with the given parameters
func NewParser(domain string, pattern []string, exclude []string, trimMarkers []string) *Parser {
	p := Parser{domain: domain, pattern: pattern, exclude: exclude, trimMarkers: trimMarkers}
	return &p
}

// GetLinks : Parse all links from the HTML document
func (p *Parser) GetLinks(htm string) ([]string, error) {
	cleanedHtm := strings.NewReader(p.trimDocument(htm))
	htmlTree, err := html.Parse(cleanedHtm)
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

	keys = p.prependDomainToLinks(p.removeExcludedLinks(p.filterPatternLinks(keys)))
	return keys, nil
}

// ExtractDocumentTitle : Extracts document title from HTML string
func (p *Parser) ExtractDocumentTitle(htm string) (string, error) {
	htmStream := strings.NewReader(p.trimDocument(htm))
	startNode, err := html.Parse(htmStream)
	title := ""

	if err != nil {
		return "", err
	}

	var crawl func(node *html.Node)
	crawl = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" {
			title = node.FirstChild.Data
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawl(child)
		}
	}

	crawl(startNode)
	return title, nil
}

func (p *Parser) filterPatternLinks(links []string) []string {
	if p.pattern == nil || len(p.pattern) == 0 {
		return links
	}

	resultLinks := make([]string, 0)

	for _, link := range links {
		for _, pattern := range p.pattern {
			if strings.HasPrefix(link, pattern) {
				resultLinks = append(resultLinks, link)
				break
			}
		}
	}

	return resultLinks
}

func (p *Parser) removeExcludedLinks(links []string) []string {
	if p.exclude == nil || len(p.exclude) == 0 {
		return links
	}

	results := make([]string, 0)
	for _, link := range links {
		shouldInclude := true

		for _, subStr := range p.exclude {
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

func (p *Parser) prependDomainToLinks(links []string) []string {
	if p.domain == "" {
		return links
	}

	newLinks := make([]string, 0)
	for _, link := range links {
		var buffer bytes.Buffer
		buffer.WriteString(p.domain)
		buffer.WriteString(link)
		newLinks = append(newLinks, buffer.String())
	}

	return newLinks
}

func (p *Parser) trimDocument(htm string) string {
	if p.trimMarkers == nil || len(p.trimMarkers) == 0 {
		return htm
	}

	parts := htm

	for _, marker := range p.trimMarkers {
		parts = strings.Split(htm, marker)[0]
	}

	return parts
}

func sanitizeLink(link string) string {
	newString := strings.ReplaceAll(link, "\"", "")
	return strings.ReplaceAll(newString, " ", "")
}
