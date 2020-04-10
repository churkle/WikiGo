package db

import (
	"database/sql"
	"fmt"
	"strconv"
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
	RetrievePageLinks(pageTitle string) []string
	RetrievePageURL(pageTitle string) string
	RetrieveAllPageTitles() []string
	RetrievePageInfo(title string) (string, bool, []string)
}

// SQLDriver : A struct that operates on the SQL db directly
type SQLDriver struct {
	db       *sql.DB
	numPages int
}

// NewSQLDriver : Creates a new SQLDriver object with the given db object
func NewSQLDriver(db *sql.DB) *SQLDriver {
	d := SQLDriver{db: db}
	d.numPages = 0
	return &d
}

// PageExists : Queries all pages in the db and finds if the page with given title exists
func (d *SQLDriver) PageExists(pageTitle string) bool {
	rs, err := d.db.Query(`SELECT title FROM pages WHERE title=$1`, pageTitle)

	if err != nil {
		fmt.Println(err)
	}

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var title string
			rs.Scan(&title)

			if pageTitle == title {
				return true
			}
		}
	}

	return false
}

// InsertPage : Inserts a new page with given title and URL into the db
func (d *SQLDriver) InsertPage(title string, url string, insertionTime time.Time) error {
	_, err := d.db.Exec(
		`INSERT INTO pages (title, url, isCrawled, lastCrawled)
		VALUES ($1, $2, $3, $4)`, title, url, "t", insertionTime.String())
	if err != nil {
		return err
	}

	return nil
}

// InsertPageTitleOnly : Inserts a page with given title into the db
func (d *SQLDriver) InsertPageTitleOnly(title string, insertionTime time.Time) error {
	_, err := d.db.Exec(
		`INSERT INTO pages (title, url, isCrawled, lastCrawled)
		VALUES ($1, $2, $3, $4)`, title, "", "f", insertionTime.String())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// InsertEdge : Inserts a new edge relationship with given source and destination IDs into db
func (d *SQLDriver) InsertEdge(srcID int, destID int) error {
	_, err := d.db.Exec(
		`INSERT INTO edges (srcID, destID)
		VALUES ($1, $2)`, srcID, destID)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePageAsCrawled : Marks the page with the given title as crawled and adds its URL
func (d *SQLDriver) UpdatePageAsCrawled(title string, url string, insertionTime time.Time) error {
	sqlStatement :=
		`UPDATE pages
	SET url = $1, isCrawled = 't'
	WHERE title = $2;`

	_, err := d.db.Exec(sqlStatement, url, title)
	if err != nil {
		return err
	}

	return nil
}

// RetrievePageLinks : Retrieves all the titles of pages that are linked to the page with the given title
func (d *SQLDriver) RetrievePageLinks(pageTitle string) []string {
	rs, err := d.db.Query(`SELECT id, title, isCrawled FROM pages WHERE title=$1`, pageTitle)
	if err != nil {
		fmt.Println(err)
	}

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var id int
			var title string
			var isCrawled string
			rs.Scan(&id, &title, &isCrawled)

			if isCrawled == "t" {
				return d.retrieveTitlesOfIDs(d.retrieveEdges(id))
			}
		}
	}

	return nil
}

// RetrievePageURL : Gets the URL of the page with the given title
func (d *SQLDriver) RetrievePageURL(pageTitle string) string {
	rs, err := d.db.Query(`SELECT title, isCrawled, url FROM pages WHERE title=$1`, pageTitle)
	if err != nil {
		fmt.Println(err)
	}

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var title string
			var url string
			var isCrawled string
			rs.Scan(&title, &url, &isCrawled)

			if pageTitle == title && isCrawled == "t" {
				return url
			}
		}
	}

	return ""
}

// RetrieveAllPageTitles : Retrieves a list of all page titles in the db
func (d *SQLDriver) RetrieveAllPageTitles() []string {
	rs, err := d.db.Query("SELECT title FROM pages")
	if err != nil {
		fmt.Println(err)
	}

	titles := make([]string, 0)

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var title string
			rs.Scan(&title)

			titles = append(titles, title)
		}
	}

	return titles
}

// RetrievePageID : Retrieves the ID of the page with the given title
func (d *SQLDriver) RetrievePageID(pageTitle string) int {
	rs, err := d.db.Query(`SELECT id FROM pages WHERE title=$1`, pageTitle)
	if err != nil {
		fmt.Println(err)
	}

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var id int
			rs.Scan(&id)

			return id
		}
	}

	return -1
}

// RetrievePageInfo : Retrieves all page info for a page given its title
func (d *SQLDriver) RetrievePageInfo(title string) (string, bool, []string) {
	rs, err := d.db.Query(`SELECT url, isCrawled FROM pages WHERE title=$1`, title)
	if err != nil {
		fmt.Println(err)
	}

	var url string
	var isCrawled bool
	var links []string

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var isCrawledString string
			rs.Scan(&url, &isCrawledString)

			if isCrawledString == "t" {
				isCrawled = true
			} else {
				isCrawled = false
			}

		}
	}

	links = d.RetrievePageLinks(title)

	return url, isCrawled, links
}

func (d *SQLDriver) retrieveEdges(srcID int) []int {
	rs, err := d.db.Query(`SELECT * FROM edges WHERE src=$1`, srcID)

	if err != nil {
		fmt.Println(err)
	}
	destIDs := make([]int, 0)

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var src int
			var dest int
			rs.Scan(&src, &dest)

			if srcID == src {
				destIDs = append(destIDs, dest)
			}
		}
	}

	return destIDs
}

func (d *SQLDriver) retrieveTitlesOfIDs(ids []int) []string {
	idListString := ""
	for index, id := range ids {
		idListString += strconv.Itoa(id)
		if index < len(ids) {
			idListString += ","
		}
	}
	rs, err := d.db.Query("SELECT id, title FROM pages WHERE id in ($1)", idListString)
	if err != nil {
		fmt.Println(err)
	}

	titles := make([]string, 0)

	if rs != nil {
		defer rs.Close()
		for rs.Next() {
			var id int
			var title string
			rs.Scan(&id, &title)

			titles = append(titles, title)
		}
	}

	return titles
}
