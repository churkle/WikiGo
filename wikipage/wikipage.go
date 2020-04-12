package wikipage

// WikiPage : a struct representing a Wiki page and its links
type WikiPage struct {
	url       string
	title     string
	links     []string
	isCrawled bool
}

// NewWikiPage : Creates a new wikipage object with the given url, title and links
func NewWikiPage(url string, title string, links []string) *WikiPage {
	return &WikiPage{url: url, title: title, links: links, isCrawled: false}
}

// NewWikiPageWithCrawlStatus : Creates a new wikipage object with the given url, title and links, and set crawled status
func NewWikiPageWithCrawlStatus(url string, title string, links []string, isCrawled bool) *WikiPage {
	return &WikiPage{url: url, title: title, links: links, isCrawled: isCrawled}
}

// NewWikiPageNoLinks : Creates a new wikipage object with no links
func NewWikiPageNoLinks(url string, title string) *WikiPage {
	return &WikiPage{url: url, title: title, links: make([]string, 0)}
}

// GetURL : Gets the wikipage's URL
func (w *WikiPage) GetURL() string {
	return w.url
}

// GetTitle : Gets the wikipage's title
func (w *WikiPage) GetTitle() string {
	return w.title
}

// GetLinks : Gets the wikipage's links
func (w *WikiPage) GetLinks() []string {
	return w.links
}

// GetCrawledStatus : Get whether the page has been crawled
func (w *WikiPage) GetCrawledStatus() bool {
	return w.isCrawled
}

// AddLink : Sets the wikipage's links
func (w *WikiPage) AddLink(link string) {
	w.links = append(w.links, link)
}
