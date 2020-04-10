package db

import (
	"WikiGo/wikipage"
	"fmt"
	"time"
)

// Service : Service wrapper that contains functionality for reading/writing cached wiki pages and their links
type Service struct {
	driver Driver
}

// NewDBService : Creates a new DB service with a driver injected
func NewDBService(driver Driver) *Service {
	return &Service{driver}
}

// AddPage : Adds a wikipage entry to the database
func (s *Service) AddPage(page *wikipage.WikiPage) error {
	var err error
	title := page.GetTitle()
	currentTime := time.Now()

	if !s.driver.PageExists(title) {
		err = s.driver.InsertPageTitleOnly(title, currentTime)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	srcID := s.driver.RetrievePageID(title)
	for _, link := range page.GetLinks() {
		if !s.driver.PageExists(link) {
			err = s.driver.InsertPageTitleOnly(link, currentTime)
			if err != nil {
				return err
			}
		}
		destID := s.driver.RetrievePageID(link)
		err = s.driver.InsertEdge(srcID, destID)
		if err != nil {
			return err
		}
	}

	s.driver.UpdatePageAsCrawled(title, page.GetURL(), currentTime)

	return nil
}

// GetPageGraph : Returns an adjacency list of a graph of wiki articles that have links to each other
func (s *Service) GetPageGraph() map[string][]string {
	titles := s.driver.RetrieveAllPageTitles()
	graph := make(map[string][]string)

	for _, title := range titles {
		links := s.driver.RetrievePageLinks(title)
		if links != nil {
			graph[title] = links
		}
	}

	return graph
}

// GetURLs : Returns a map of page titles and their URLs
func (s *Service) GetURLs() map[string]string {
	titles := s.driver.RetrieveAllPageTitles()
	urlMap := make(map[string]string)

	for _, title := range titles {
		url := s.driver.RetrievePageURL(title)
		if url != "" {
			urlMap[title] = url
		}
	}

	return urlMap
}

// GetPage : Returns a wikipage object of a page title if it exists in the db
func (s *Service) GetPage(title string) *wikipage.WikiPage {
	url, isCrawled, links := s.driver.RetrievePageInfo(title)
	return wikipage.NewWikiPageWithCrawlStatus(url, title, links, isCrawled)
}
