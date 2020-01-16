package db

import (
	"WikiGo/wikipage"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type TestDatabaseDriver struct {
	db       *sql.DB
	numPages int
}

func NewTestDBDriver(db *sql.DB) *TestDatabaseDriver {
	d := TestDatabaseDriver{db: db}
	d.numPages = 0
	return &d
}

func (d *TestDatabaseDriver) PageExists(pageTitle string) bool {
	rs, _ := d.db.Query("SELECT pages")
	defer rs.Close()

	for rs.Next() {
		var id int
		var title string
		var url string
		var isCrawled string
		var lastCrawled string
		rs.Scan(&id, &title, &url, &isCrawled, &lastCrawled)

		if pageTitle == title {
			return true
		}
	}

	return false
}

func (d *TestDatabaseDriver) InsertPage(title string, url string, insertionTime time.Time) error {
	_, err := d.db.Exec(
		`INSERT INTO pages (title, url, isCrawled, lastCrawled)
		VALUES ($1, $2, $3, $4)`, title, url, "t", insertionTime.String())
	if err != nil {
		return err
	}

	return nil
}

func (d *TestDatabaseDriver) InsertEdge(srcID int, destID int) error {
	_, err := d.db.Exec(
		`INSERT INTO edges (srcID, destID)
		VALUES ($1, $2)`, srcID, destID)
	if err != nil {
		return err
	}

	return nil
}

func (d *TestDatabaseDriver) UpdatePageAsCrawled(title string, url string, insertionTime time.Time) error {
	_, err := d.db.Exec(
		`UPDATE pages
		SET url = $1, isCrawled = "t", lastCrawled = $2
		WHERE title = $3`, url, insertionTime.String(), title)
	if err != nil {
		return err
	}

	return nil
}

func (d *TestDatabaseDriver) RetrievePageURLAndLinks(pageTitle string) (string, []string) {
	rs, _ := d.db.Query("SELECT pages")
	defer rs.Close()

	for rs.Next() {
		var id int
		var title string
		var url string
		var isCrawled string
		var lastCrawled string
		rs.Scan(&id, &title, &url, &isCrawled, &lastCrawled)

		if pageTitle == title && isCrawled == "t" {
			return url, d.retrieveTitles(d.retrieveEdges(id))
		}
	}

	return "", nil
}

func (d *TestDatabaseDriver) RetrievePageID(pageTitle string) int {
	rs, _ := d.db.Query("SELECT pages")
	defer rs.Close()

	for rs.Next() {
		var id int
		var title string
		var url string
		var isCrawled string
		var lastCrawled string
		rs.Scan(&id, &title, &url, &isCrawled, &lastCrawled)

		if pageTitle == title {
			return id
		}
	}

	return -1
}

func (d *TestDatabaseDriver) retrieveEdges(srcID int) []int {
	rs, _ := d.db.Query("SELECT edges")
	defer rs.Close()

	destIDs := make([]int, 0)
	for rs.Next() {
		var src int
		var dest int
		rs.Scan(&src, &dest)

		if srcID == src {
			destIDs = append(destIDs, src)
		}
	}

	return destIDs
}

func (d *TestDatabaseDriver) retrieveTitles(ids []int) []string {
	rs, _ := d.db.Query("SELECT pages")
	defer rs.Close()

	titles := make([]string, 0)
	for rs.Next() {
		var id int
		var title string
		var url string
		var isCrawled string
		var lastCrawled string
		rs.Scan(&id, &title, &url, &isCrawled, &lastCrawled)

		for _, x := range ids {
			if id == x {
				titles = append(titles, title)
			}
		}
	}

	return titles
}

func assertSameWikiPage(t *testing.T, expected *wikipage.WikiPage, result *wikipage.WikiPage) {
	t.Helper()

	if (expected.GetURL() != result.GetURL()) || (expected.GetTitle() != result.GetTitle()) {
		t.Errorf("Expected '%q' but got '%q'", expected, result)
	}

	for index, expectedVal := range expected.GetLinks() {
		if result.GetLinks()[index] != expectedVal {
			t.Errorf("Expected '%q' but got '%q'", expectedVal, result.GetLinks()[index])
		}
	}
}

func TestInsertPage(t *testing.T) {
	t.Run("Test inserting one page", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		pageRows := sqlmock.NewRows([]string{"id", "title", "url", "isCrawled", "lastCrawled"})
		edgeRows := sqlmock.NewRows([]string{"srcID", "destID"})
		mock.ExpectQuery("SELECT pages").WillReturnRows(pageRows)
		mock.ExpectQuery("SELECT edges").WillReturnRows(edgeRows)
		mock.ExpectExec("INSERT INTO pages").WithArgs("Example Page", "www.example.com", "t", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO edges").WithArgs(0, 1).WillReturnResult(sqlmock.NewResult(1, 1))

		testDriver := NewTestDBDriver(db)

		testDBService := NewDBService(testDriver)
		testPage := wikipage.NewWikiPage("www.example.com", "Example Page", []string{"Example 1"})
		testDBService.AddPage(testPage)

		expected := map[string][]string{"Example Page": []string{"Example 1"}}
		result := testDBService.GetPageGraph()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected '%q' but got '%q'", expected, result)
		}
	})
}
