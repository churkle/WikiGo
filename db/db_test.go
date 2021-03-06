package db

import (
	"WikiGo/wikipage"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func assertSameWikiPage(t *testing.T, expected *wikipage.WikiPage, result *wikipage.WikiPage) {
	t.Helper()

	if (expected.GetURL() != result.GetURL()) || (expected.GetTitle() != result.GetTitle()) {
		t.Errorf("Expected '%q' but got '%q'", expected.GetURL()+expected.GetTitle(), result.GetURL()+result.GetTitle())
	}

	for index, expectedVal := range expected.GetLinks() {
		if result.GetLinks()[index] != expectedVal {
			t.Errorf("Expected '%q' but got '%q'", expectedVal, result.GetLinks()[index])
		}
	}
}

type TestPage struct {
	title     string
	url       string
	isCrawled string
}

func TestInsertPage(t *testing.T) {
	t.Run("Test inserting one page", func(t *testing.T) {
		testObject := TestPage{title: "Example Page", url: "www.example.com", isCrawled: "t"}
		testLink := TestPage{title: "Example 1", url: "", isCrawled: "f"}

		db, mock, err := sqlmock.New()
		if err != nil {
			fmt.Println("failed to open sqlmock database:", err)
		}
		defer db.Close()

		mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(sqlmock.NewRows([]string{"title"}))
		mock.ExpectExec("INSERT INTO pages").WithArgs(testObject.title, "", "f", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(0))
		mock.ExpectQuery(`SELECT`).WithArgs(testLink.title).WillReturnRows(
			sqlmock.NewRows([]string{"title"}))
		mock.ExpectExec("INSERT INTO pages").WithArgs(testLink.title, testLink.url, testLink.isCrawled, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectQuery(`SELECT`).WithArgs(testLink.title).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(1))
		mock.ExpectExec("INSERT INTO edges").WithArgs(0, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE pages").WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`SELECT`).WillReturnRows(
			sqlmock.NewRows([]string{"title"}).
				AddRow(testObject.title).
				AddRow(testLink.title))
		mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "isCrawled"}).
				AddRow(0, testObject.title, testObject.isCrawled))
		mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{"srcID", "destID"}).AddRow(0, 1))
		mock.ExpectQuery(`SELECT`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "title"}).
				AddRow(1, testLink.title))
		mock.ExpectQuery(`SELECT`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "url", "isCrawled", "lastCrawled"}).
				AddRow(0, testObject.title, testObject.url, testObject.isCrawled, "NA").
				AddRow(1, testLink.title, testLink.url, testLink.isCrawled, "NA"))

		testDriver := NewSQLDriver(db)
		testDBService := NewDBService(testDriver)
		testPage := wikipage.NewWikiPage(testObject.url, testObject.title, []string{"Example 1"})
		testDBService.AddPage(testPage)

		expected := map[string][]string{testObject.title: []string{"Example 1"}}
		result := testDBService.GetPageGraph()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Expected '%q' but got '%q'", expected, result)
		}
	})
}

func TestGetPage(t *testing.T) {
	testObject := TestPage{title: "Example Page", url: "www.example.com", isCrawled: "t"}
	testLink := TestPage{title: "Example 1", url: "", isCrawled: "f"}

	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(sqlmock.NewRows([]string{"title"}))
	mock.ExpectExec("INSERT INTO pages").WithArgs(testObject.title, "", "f", sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).
			AddRow(0))
	mock.ExpectQuery(`SELECT`).WithArgs(testLink.title).WillReturnRows(
		sqlmock.NewRows([]string{"title"}))
	mock.ExpectExec("INSERT INTO pages").WithArgs(testLink.title, testLink.url, testLink.isCrawled, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT`).WithArgs(testLink.title).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	mock.ExpectExec("INSERT INTO edges").WithArgs(0, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE pages").WillReturnResult(sqlmock.NewResult(1, 1))

	testDriver := NewSQLDriver(db)
	testDBService := NewDBService(testDriver)
	testPage := wikipage.NewWikiPageWithCrawlStatus(testObject.url, testObject.title, []string{"Example 1"}, true)
	testDBService.AddPage(testPage)

	mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(sqlmock.NewRows([]string{"url", "isCrawled"}).
		AddRow(testObject.url, testObject.isCrawled))
	mock.ExpectQuery(`SELECT`).WithArgs(testObject.title).WillReturnRows(
		sqlmock.NewRows([]string{"id", "title", "isCrawled"}).
			AddRow(0, testObject.title, testObject.isCrawled))
	mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{"srcID", "destID"}).AddRow(0, 1))
	mock.ExpectQuery(`SELECT`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "title"}).
			AddRow(1, testLink.title))

	resultPage := testDBService.GetPage(testObject.title)
	assertSameWikiPage(t, testPage, resultPage)
}
