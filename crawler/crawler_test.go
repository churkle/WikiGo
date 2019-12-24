package crawler

import (
	"testing"
)

func TestGetHTMLReaderFromURL(t *testing.T) {
	t.Run("Retrieving HTML reader from a URL", func(t *testing.T) {
		result, err := GetHTMLReaderFromURL("https://www.wikipedia.org/")

		if err != nil {
			t.Error(err)
		}

		if result == "" {
			t.Error("Couldn't retrieve URL")
		}
	})
}
