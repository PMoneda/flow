package gonnie

import (
	"bufio"
	"bytes"
	"strings"
)

func msg(next func(), e *ExchangeMessage, out Message, u uri, params ...interface{}) error {
	if len(params) == 0 {
		panic("Message Letter required")
	}
	letter := strings.NewReader(params[0].(string))
	buf := bufio.NewReader(letter)
	line, _, err := buf.ReadLine()
	isHeader := false
	isBody := false
	buff := bytes.Buffer{}
	for err == nil {
		str := strings.TrimSpace(string(line))
		if str == "Header:" {
			isHeader = true
		} else if str == "Body:" {
			isHeader = false
			isBody = true
		} else if isHeader {
			setHeaders(e, str)
		} else if isBody {
			buff.Write(line)
		}
		line, _, err = buf.ReadLine()
	}
	e.SetBody(buff.String())
	out <- e
	next()
	return nil
}
func setHeaders(e *ExchangeMessage, line string) {
	if line == "" {
		return
	}
	kv := strings.Split(line, ":")
	if len(kv) < 2 {
		panic("Invalid Header: " + line)
	}
	e.SetHeader(kv[0], kv[1])

}
