package parser

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestReadHTML(t *testing.T) {
	body, _ := ioutil.ReadFile("test.html")
	result := string(body)

	expected := `<html>
<head>
<title>Test website</title>
</head>
<body>
<p>Test paragraph</p>
<p>Here's a <a href=http://www.yahoo.com/>Link!</a></p>
</body>
</html>`

	if result != expected {
		t.Errorf("Expected '%q' but got '%q'", expected, result)
	}
}

func TestGetLinks(t *testing.T) {
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

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%q' but got '%q'", expected, result)
	}
}
