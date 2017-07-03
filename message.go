package flow

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"strings"
)

func messageConnector(e *ExchangeMessage, u URI, params ...interface{}) error {

	if len(params) < 2 {
		return errors.New("Message Letter required")
	}
	var fncs template.FuncMap
	var body interface{}
	var tmpl string
	if u.GetOption("source") == "inline" {
		body = params[0]
		tmpl = params[1].(string)
		if len(params) > 2 {
			fncs = params[2].(template.FuncMap)
		}
	} else {
		body = e.body
	}
	parsed, errTmpl := parseTemplate(body, tmpl, fncs)

	if errTmpl != nil {
		return errTmpl
	}
	letter := strings.NewReader(parsed)
	buf := bufio.NewReader(letter)
	line, _, err := buf.ReadLine()
	buff := bytes.Buffer{}
	isXML := false
	for err == nil {
		str := strings.TrimSpace(string(line))
		if strings.HasPrefix(str, "##") {
			str = strings.Replace(str, "##", "", 1)
			err := setHeaders(e, str)
			if strings.Contains(str, "text/xml") {
				buff.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
				buff.WriteString("\n")
				isXML = true
			}
			if err != nil {
				return err
			}
		} else if str != "" {
			buff.Write(line)
			if isXML {
				buff.WriteString("\n")
			}
		}
		line, _, err = buf.ReadLine()
	}
	s := buff.String()
	e.SetBody(s)
	return nil
}
func setHeaders(e *ExchangeMessage, line string) error {
	if line == "" {
		return nil
	}
	line = strings.TrimSpace(line)
	kv := strings.Split(line, ":")
	if len(kv) < 2 {
		return errors.New("Invalid Header: " + line)
	}
	e.SetHeader(kv[0], kv[1])
	return nil
}
