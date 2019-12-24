package parser

import (
	"io/ioutil"
	"strings"
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

func TestGetLinks(t *testing.T) {
	t.Run("Using an empty body", func(t *testing.T) {
		testBody := ""

		result, err := GetLinks(strings.NewReader(testBody))

		if err != nil {
			t.Error(err)
		}

		expected := []string{}

		assertSameSlice(t, result, expected)
	})

	t.Run("Body with no links", func(t *testing.T) {
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Test paragraph</p>
</body>
</html>`

		result, err := GetLinks(strings.NewReader(testBody))

		if err != nil {
			t.Error(err)
		}

		expected := []string{}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding a single link", func(t *testing.T) {
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Test paragraph</p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
</body>
</html>`

		result, err := GetLinks(strings.NewReader(testBody))

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding multiple links", func(t *testing.T) {
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
<p>Test paragraph</p>
<p>Here's another <a href=http://www.google.com/>Link!</a></p>
</body>
</html>`

		result, err := GetLinks(strings.NewReader(testBody))

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.google.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Don't add duplicates", func(t *testing.T) {
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Here's another <a href=http://www.google.com/>Link!</a></p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
<p>Test paragraph</p>
<p>Here's another <a href=http://www.google.com/>Link!</a></p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
</body>
</html>`

		result, err := GetLinks(strings.NewReader(testBody))

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.google.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding nested links", func(t *testing.T) {
		body, _ := ioutil.ReadFile("test.html")
		htm := string(body)
		result, err := GetLinks(strings.NewReader(htm))

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		assertSameSlice(t, result, expected)
	})
}

func TestFilterPatternLinks(t *testing.T) {
	t.Run("Filter empty pattern", func(t *testing.T) {
		patterns := []string{""}
		testLinks := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		expected := []string{"http://www.yahoo.com/", "http://www.apple.com/"}

		assertSameSlice(t, FilterPatternLinks(testLinks, patterns), expected)
	})

	t.Run("Filter http pattern", func(t *testing.T) {
		patterns := []string{"http://"}
		testLinks := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		expected := []string{"http://www.yahoo.com/", "http://www.apple.com/"}

		assertSameSlice(t, FilterPatternLinks(testLinks, patterns), expected)
	})

	t.Run("Filter pattern that fits none", func(t *testing.T) {
		patterns := []string{"https://"}
		testLinks := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		expected := []string{}

		assertSameSlice(t, FilterPatternLinks(testLinks, patterns), expected)
	})

	t.Run("Filter specific domain", func(t *testing.T) {
		patterns := []string{"http://www.yahoo.com"}
		testLinks := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		expected := []string{"http://www.yahoo.com/"}

		assertSameSlice(t, FilterPatternLinks(testLinks, patterns), expected)
	})
}

func TestRemoveExcludedLinks(t *testing.T) {
	t.Run("Exclude specific domain", func(t *testing.T) {
		exclude := []string{"apple"}
		testLinks := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		expected := []string{"http://www.yahoo.com/"}

		assertSameSlice(t, RemoveExcludedLinks(testLinks, exclude), expected)
	})
}

func TestPrependDomain(t *testing.T) {
	t.Run("Prepend empty string", func(t *testing.T) {
		prefix := ""
		testLinks := []string{"www.yahoo.com/", "www.apple.com/"}
		expected := []string{"www.yahoo.com/", "www.apple.com/"}

		assertSameSlice(t, PrependDomainToLinks(testLinks, prefix), expected)
	})

	t.Run("Prepend http", func(t *testing.T) {
		prefix := "http://"
		testLinks := []string{"www.yahoo.com/", "www.apple.com/"}
		expected := []string{"http://www.yahoo.com/", "http://www.apple.com/"}

		assertSameSlice(t, PrependDomainToLinks(testLinks, prefix), expected)
	})
}

func TestTrimDocument(t *testing.T) {
	t.Run("Trim html", func(t *testing.T) {
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Here's another <a href=http://www.google.com/>Link!</a></p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
<p>Test paragraph</p>
<p>Here's another <a href=http://www.google.com/>Link!</a></p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
</body>
</html>`
		expected := `<html>
<head>
<title>Test website</title>
</head>
`

		result := TrimDocument(testBody, "<body>")
		if result != expected {
			t.Errorf("Expected '%q' but got '%q'", expected, result)
		}
	})
}
