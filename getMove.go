package main

import (
	"github.com/PuerkitoBio/goquery"
)

func getMoveFromHTML(doc *goquery.Document) (MoveResult, error) {
	ReturnMove := MoveResult{}
	doc.Find("table.infobox").Find("tr").Each(func(i int, table *goquery.Selection) {
		if i == 0 {
			ReturnMove.Name = table.Text()
		}
		if i == 2 || i == 3 {
			table.Find("tr").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					ReturnMove.Type = s.Find("td").Text()
				case 1:
					ReturnMove.Category = s.Find("td").Text()
				case 2:
					ReturnMove.PP = s.Find("td").Text()
				case 3:
					ReturnMove.Power = s.Find("td").Text()
				case 4:
					ReturnMove.Accuracy = s.Find("td").Text()
				}
			})
		}
	})
	return ReturnMove, nil
}

/*func (cfg *Config) getNamedURLsFromHTML(htmlBody string) ([]string, error) {

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
} */
