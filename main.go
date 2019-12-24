package main

import (
	"WikiGo/crawler"
	"WikiGo/parser"
	"fmt"
	"strings"
)

func main() {
	htm, _ := crawler.GetHTMLReaderFromURL("https://en.wikipedia.org/wiki/UK_miners'_strike_(1984%E2%80%9385)")
	fmt.Println(parser.GetLinks(strings.NewReader(htm)))
}
