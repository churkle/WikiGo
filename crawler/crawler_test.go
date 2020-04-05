package crawler

import (
	"testing"
)

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
			"", nil, nil, nil, 3, false, nil)

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
			"", nil, nil, nil, 3, false, nil)

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
			"", nil, nil, nil, 2, false, nil)

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
			"", nil, nil, nil, 4, false, nil)

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
			"", nil, nil, nil, 5, false, nil)

		expected := []string{"CyclePage", "CyclePage2", "Page 1", "Page 2", "Page 3", "Page 4"}
		path, err := myCrawler.GetShortestPathToArticle()

		if err != nil {
			t.Error(err)
		}

		assertSameSlice(t, path, expected)
	})
}
