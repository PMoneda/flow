package flow

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//crawler://google.com?schema=https"#"
func crawlerConnector(next func(), e *ExchangeMessage, out Message, u URI, params ...interface{}) error {
	schema := "http"
	if len(params) == 0 {
		out <- e
		next()
		return errors.New("Wrong argument, you need to pass document query")
	}
	if u.GetHost() != "" {
		if len(params) == 2 {
			schema = params[1].(string)
		}
	}
	url := strings.Replace(u.GetRaw(), "crawler", schema, -1)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		e.SetHeader("status", "500")
		e.SetHeader("error", err.Error())
		out <- e
		next()
		return err
	}
	e.SetHeader("status", "200")
	selection := doc.Find(params[0].(string))
	find := make([]string, selection.Size(), selection.Size())
	selection.Each(func(i int, s *goquery.Selection) {
		find[i] = s.Text()
	})
	e.SetBody(find)
	out <- e
	next()
	return nil
}
