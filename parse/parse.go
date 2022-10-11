package parse

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

func ParseData(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("tr > td.vline > a:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			g.Get(r.JoinURL(href), func(g *geziyor.Geziyor, r *client.Response) {
				name := r.HTMLDoc.Find(".col-content > h1:first-child").Text()
				author := r.HTMLDoc.Find(".table-bookinfo span").Text()
				poem := r.HTMLDoc.Find(".table-quote .quote-t > span").Text()

				g.Exports <- map[string]interface{}{
					"Название": strings.TrimSpace(name),
					"Автор":    strings.TrimSpace(author),
					"Текст":    strings.TrimSpace(poem),
				}
			})
		}
	})
}
