package main

import (
	"WikiGo/crawler"
	"fmt"
)

func main() {
	patterns := []string{"/wiki/"}
	exclude := []string{"Wikipedia:", "Special:", "Help:", "Books:", "File:", ".jpg"}
	myCrawler := crawler.NewCrawler("https://en.wikipedia.org/wiki/UK_miners'_strike_(1984%E2%80%9385)",
		"https://en.wikipedia.org/wiki/Lawrence_Daly",
		"https://en.wikipedia.org", patterns, exclude, ">Notes<", 4)

	path, err := myCrawler.GetShortestPathToArticle()

	if err != nil {
		fmt.Print(err)
	}

	for _, site := range path {
		fmt.Print(site + " -> ")
	}
}
