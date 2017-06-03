package flow

import (
	"bytes"
	"html/template"
)

func templateConnector(next func(), m *ExchangeMessage, out Message, u URI, params ...interface{}) error {
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

func parseTemplate(obj interface{}, _tmpl string, fncs template.FuncMap) (string, error) {
	buf := bytes.NewBuffer(nil)
	hash := sha(_tmpl)
	t, ok := tmpl[hash]
	if !ok {
		ltmpl.Lock()
		t = template.Must(template.New(hash).Funcs(fncs).Parse(_tmpl))
		tmpl[hash] = t
		ltmpl.Unlock()
	}
	err := t.ExecuteTemplate(buf, hash, obj)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
