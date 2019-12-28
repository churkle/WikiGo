package main

import (
	"WikiGo/crawler"
	"fmt"
)

func main() {
	patterns := []string{"/wiki/"}
	exclude := []string{"Wikipedia:", "Special:", "Help:", "Books:", "File:", ".jpg"}
	trimMarkers := []string{">Notes<", ">References<", ">See also<", `#External_links">`, `id="catlinks"`}
	myCrawler := crawler.NewCrawler("https://en.wikipedia.org/wiki/UK_miners'_strike_(1984%E2%80%9385)",
		"https://en.wikipedia.org/wiki/Berkhamsted",
		"https://en.wikipedia.org", patterns, exclude, trimMarkers, 3)

	path, err := myCrawler.GetShortestPathToArticle()

	if err != nil {
		fmt.Print(err)
	}

	if path == nil {
		fmt.Println("FAILED")
	}
}
