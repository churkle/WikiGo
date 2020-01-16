package db

import (
	"time"
)

// Driver : interface that has the basic operations for data store interaction
type Driver interface {
	PageExists(pageTitle string) bool
	RetrievePageID(pageTitle string) int
	InsertPage(title string, url string, insertionTime time.Time) error
	InsertPageTitleOnly(title string, insertionTime time.Time) error
	UpdatePageAsCrawled(title string, url string, insertionTime time.Time) error
	InsertEdge(sourceID int, destID int) error
	RetrievePageURLAndLinks(pageTitle string) (string, []string)
}
