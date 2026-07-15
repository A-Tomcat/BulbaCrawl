package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HtmlToDoc(htmlText string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, err
	}
	return doc, nil
}
func (cfg *Config) getNamedURLsFromHTML(doc *goquery.Document) ([]string, error) {
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
			if strings.Contains(parsedHref.Path, cfg.SearchName) {
				resolvedURL := cfg.BaseURL.ResolveReference(parsedHref)
				links = append(links, resolvedURL.String())
			}
		}
	})
	if crawlErr != nil {
		return nil, crawlErr
	}
	return links, nil
} //Gets URL strings of all the TCG Cards of the Searched Pokemon

func (cfg *Config) getHTML() (string, error) {
	name_parts := strings.Split(cfg.SearchName, " ")
	searchname_url := strings.Join(name_parts, "_")
	SearchPath, err := url.Parse("wiki/" + searchname_url + "_(" + cfg.Category + ")")
	if err != nil {
		return "", err
	}
	fmt.Println(SearchPath)
	SearchURL := cfg.BaseURL.ResolveReference(SearchPath)
	req, err := http.NewRequest("GET", SearchURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "BulbaCrawler/1.0")
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error, Status Code: %d", res.StatusCode)
	}
	content_type := res.Header.Get("Content-Type")
	if !strings.Contains(content_type, "text/html") {
		return "", fmt.Errorf("unexpected content type: %s", content_type)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	htmlString := string(bytes)
	return htmlString, nil
}

func (cfg *Config) CaseTCG() error {
	html, err := cfg.getHTML()
	if err != nil {
		return err
	}
	doc, err := HtmlToDoc(html)
	if err != nil {
		return err
	}
	urls, err := cfg.getNamedURLsFromHTML(doc)
	if err != nil {
		return err
	}
	cfg.result.TCG.NumberOfCards = len(urls)
	return nil
}
