package parser

import (
	"io/ioutil"
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
		p := NewParser("", []string{}, []string{}, "")

		result, err := p.GetLinks("")

		if err != nil {
			t.Error(err)
		}

		expected := []string{}

		assertSameSlice(t, result, expected)
	})

	t.Run("Body with no links", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")

		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Test paragraph</p>
</body>
</html>`

		result, err := p.GetLinks(testBody)

		if err != nil {
			t.Error(err)
		}

		expected := []string{}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding a single link", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")
		testBody := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Test paragraph</p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
</body>
</html>`

		result, err := p.GetLinks(testBody)

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding multiple links", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")
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

		result, err := p.GetLinks(testBody)

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.google.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Don't add duplicates", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")
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

		result, err := p.GetLinks(testBody)

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.google.com/"}

		assertSameSlice(t, result, expected)
	})

	t.Run("Finding nested links", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")
		body, _ := ioutil.ReadFile("test.html")
		htm := string(body)
		result, err := p.GetLinks(htm)

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/", "http://www.apple.com/"}
		assertSameSlice(t, result, expected)
	})

	t.Run("Find links after trim", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "<ul>")
		body, _ := ioutil.ReadFile("test.html")
		htm := string(body)
		result, err := p.GetLinks(htm)

		if err != nil {
			t.Error(err)
		}

		expected := []string{"http://www.yahoo.com/"}
		assertSameSlice(t, result, expected)
	})
}

func TestExtractDocumentTitle(t *testing.T) {
	t.Run("Find title of document", func(t *testing.T) {
		p := NewParser("", []string{}, []string{}, "")
		body, _ := ioutil.ReadFile("test.html")
		htm := string(body)
		result, err := p.ExtractDocumentTitle(htm)

		if err != nil {
			t.Error(err)
		}

		expected := "Test website"
		if result != expected {
			t.Errorf("Expected '%q' but got '%q'", expected, result)
		}
	})
}
