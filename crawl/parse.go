package crawl

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
)
// Creates and crawls a document from a given url
func crawl(in string) CrawlerResult {
	url, err := normalizeUrl(in)
	if err != nil {
		return CrawlerResult{Error: err}
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return CrawlerResult{Error: err}
	}
	links := findLinks(doc)
	return CrawlerResult{Url: url, Data: links}
}

// Extracts links from a document
func findLinks(doc *goquery.Document) []string {
	links := make([]string, 0)
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		links = append(links, link)
	})
	return links
}

func normalizeUrl(in string) (out string, err error) {
	url, err := url.Parse(in)
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	out = url.String()
	return
}
