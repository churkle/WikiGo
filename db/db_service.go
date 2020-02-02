package db

import (
	"WikiGo/wikipage"
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

	if !s.driver.PageExists(page.GetTitle()) {
		err = s.driver.InsertPageTitleOnly(page.GetTitle(), time.Now())
		if err != nil {
			return err
		}
	}

	srcID := s.driver.RetrievePageID(page.GetTitle())
	for _, link := range page.GetLinks() {
		if !s.driver.PageExists(link) {
			err = s.driver.InsertPageTitleOnly(link, time.Now())
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

	s.driver.UpdatePageAsCrawled(page.GetTitle(), page.GetURL(), time.Now())

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
//func (s *Service) GetURLs() map[string]string {

//}
