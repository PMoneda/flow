package gonnie

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func getClient(skip string) *http.Client {
	_skip := false
	if skip != "" {
		_skip, _ = strconv.ParseBool(skip)
	}
	cfg := &tls.Config{
		InsecureSkipVerify: _skip,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg,
		},
	}
	return client
}

func httpConector(next func(), e *ExchangeMessage, out Message, u uri, params ...interface{}) error {
	skip := u.options.Get("insecureSkipVerify")
	client := getClient(skip)
	var req *http.Request
	var err error
	b := e.GetBody()
	var body io.Reader
	switch t := b.(type) {
	case string:
		body = strings.NewReader(t)
	default:
		if j, err := json.Marshal(b); err != nil {
			panic(err)
		} else {
			body = strings.NewReader(string(j))
		}
	}
	if len(u.options) > 0 {
		req, err = http.NewRequest(u.options.Get("method"), u.options.Get("url"), body)
	} else {
		req, err = http.NewRequest("GET", u.raw, body)
	}
	if err != nil {
		return err
	}
	authMethod := u.options.Get("auth")
	if authMethod == "basic" {
		req.SetBasicAuth(u.options.Get("user"), u.options.Get("password"))
	}
	header := e.head
	keys := header.ListKeys()
	for _, key := range keys {
		req.Header.Add(key, header.Get(key))
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return errResponse
	}
	newData := NewExchangeMessage()
	for k := range resp.Header {
		newData.SetHead(k, resp.Header.Get(k))
	}
	newData.SetBody(string(data))
	out <- newData
	next()
	return nil
}
