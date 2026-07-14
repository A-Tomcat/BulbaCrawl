package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}
	var links []string
	var crawlErr error
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			parsedHref, err := url.Parse(href)
			if err != nil {
				crawlErr = err
				return
			}
			resolvedURL := baseURL.ResolveReference(parsedHref)
			links = append(links, resolvedURL.String())
		}
	})
	if crawlErr != nil {
		return nil, crawlErr
	}
	return links, nil
}
