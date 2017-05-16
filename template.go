package gonnie

import (
	"bytes"
	"html/template"
)

func templateConector(next func(), m *ExchangeMessage, out Message, u Uri, params ...interface{}) error {
	buf := bytes.NewBuffer(nil)
	if len(params) > 1 && params[1] != nil {
		fncs := params[1].(template.FuncMap)
		for k, v := range fncs {
			funcMap[k] = v
		}
	}
	t := template.Must(template.New("transform").Funcs(funcMap).Parse(params[0].(string)))
	err := t.ExecuteTemplate(buf, "transform", m.body)
	if err != nil {
		panic(err)
	}
	m.body = buf.String()
	out <- m
	return nil
}
