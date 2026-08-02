package main

import (
	"github.com/PuerkitoBio/goquery"
)

func getMoveFromDoc(doc *goquery.Document) (MoveResult, error) {
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
