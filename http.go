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
func httpConector(next func(), e *ExchangeMessage, out Message, u Uri, params ...interface{}) error {
	newData := NewExchangeMessage()
	var skip string
	var opts map[string]string
	method := "GET"
	authMethod := ""
	username := ""
	password := ""
	if len(params) > 0 {
		opts = params[0].(map[string]string)
		skip = opts["insecureSkipVerify"]
		method = opts["method"]
		authMethod = opts["auth"]
		username = opts["username"]
		password = opts["password"]
	}
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
	req, err = http.NewRequest(method, u.raw, body)
	if err != nil {
		newData.SetBody(err)
		out <- newData
		return err
	}
	if authMethod == "basic" {
		req.SetBasicAuth(username, password)
	}
	header := e.head
	keys := header.ListKeys()
	for _, key := range keys {
		req.Header.Add(key, header.Get(key))
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		newData.SetBody(errResp)
		out <- newData
		return errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		newData.SetBody(errResponse)
		out <- newData
		return errResponse
	}

	for k := range resp.Header {
		newData.SetHeader(k, resp.Header.Get(k))
	}
	newData.SetBody(string(data))
	out <- newData
	next()
	return nil
}
