package gonnie

import (
	"bytes"
	"html/template"
)

func templateConector(next func(), m *ExchangeMessage, out Message, u Uri, params ...interface{}) error {
	var fncs template.FuncMap
	if len(params) > 1 && params[1] != nil {
		fncs = params[1].(template.FuncMap)
	}
	s, err := parseTemplate(m.body, params[0].(string), fncs)
	if err != nil {
		return err
	}
	m.body = s
	out <- m
	return nil
}

func parseTemplate(obj interface{}, tmpl string, fncs template.FuncMap) (string, error) {
	buf := bytes.NewBuffer(nil)
	for k, v := range fncs {
		funcMap[k] = v
	}
	t := template.Must(template.New("transform").Funcs(funcMap).Parse(tmpl))
	err := t.ExecuteTemplate(buf, "transform", obj)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
