package crawler

import (
	"WikiGo/db"
	"testing"
	"time"
)

// TestDBDriver : A test DB driver that doesn't actually do anything
type TestDBDriver struct {
}

func (td *TestDBDriver) PageExists(pageTitle string) bool {
	return false
}

func (td *TestDBDriver) RetrievePageID(pageTitle string) int {
	return -1
}

func (td *TestDBDriver) InsertPage(title string, url string, insertionTime time.Time) error {
	return nil
}

func (td *TestDBDriver) InsertPageTitleOnly(title string, insertionTime time.Time) error {
	return nil
}

func (td *TestDBDriver) UpdatePageAsCrawled(title string, url string, insertionTime time.Time) error {
	return nil
}

func (td *TestDBDriver) InsertEdge(sourceID int, destID int) error {
	return nil
}

func (td *TestDBDriver) RetrievePageLinks(pageTitle string) []string {
	return nil
}

func (td *TestDBDriver) RetrievePageURL(pageTitle string) string {
	return ""
}

func (td *TestDBDriver) RetrieveAllPageTitles() []string {
	return nil
}

func (td *TestDBDriver) RetrievePageInfo(title string) (string, bool, []string) {
	return "", false, nil
}

func assertSameSlice(t *testing.T, result, expected []string) {
	t.Helper()

	if len(result) != len(expected) {
		t.Errorf("Expected '%q' but got '%q'", expected, result)
	}

	for _, val1 := range result {
		exists := false
		for _, val2 := range expected {
			if val1 == val2 {
				exists = true
			}
		}

		if !exists {
			t.Errorf("Expected '%q' but got '%q'", expected, result)
		}
	}
}

func TestGetHTMLReaderFromURL(t *testing.T) {
	t.Run("Simple crawl with just 2 pages", func(t *testing.T) {
		myCrawler := NewCrawler(`./testHTML/page1.html`,
			`./testHTML/page2.html`,
			"", nil, nil, nil, 3, false, db.NewDBService(&TestDBDriver{}))

		expected := []string{"Page 1", "Page 2"}
		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		assertSameSlice(t, path, expected)
	})

	t.Run("Test 3 pages", func(t *testing.T) {
		myCrawler := NewCrawler(`./testHTML/page1.html`,
			`./testHTML/page3.html`,
			"", nil, nil, nil, 3, false, db.NewDBService(&TestDBDriver{}))

		expected := []string{"Page 1", "Page 2", "Page 3"}
		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		assertSameSlice(t, path, expected)
	})

	t.Run("Test depth limit", func(t *testing.T) {
		myCrawler := NewCrawler(`./testHTML/page1.html`,
			`./testHTML/page4.html`,
			"", nil, nil, nil, 2, false, db.NewDBService(&TestDBDriver{}))

		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		if path != nil {
			t.Error("Test Failed: Should not return a path since depth is too low")
		}
	})

	t.Run("Only report shortest path found", func(t *testing.T) {
		myCrawler := NewCrawler(`./testHTML/connectPage.html`,
			`./testHTML/page4.html`,
			"", nil, nil, nil, 4, false, db.NewDBService(&TestDBDriver{}))

		expected := []string{"ConnectPage", "Page 4"}
		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		assertSameSlice(t, path, expected)
	})

	t.Run("Only report shortest path found", func(t *testing.T) {
		myCrawler := NewCrawler(`./testHTML/cyclePage.html`,
			`./testHTML/page4.html`,
			"", nil, nil, nil, 5, false, db.NewDBService(&TestDBDriver{}))

		expected := []string{"CyclePage", "CyclePage2", "Page 1", "Page 2", "Page 3", "Page 4"}
		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		assertSameSlice(t, path, expected)
	})
}
