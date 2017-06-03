package flow

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var cfg *tls.Config = &tls.Config{
	InsecureSkipVerify: true,
}
var client *http.Client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:     cfg,
		Dial:                defaultDialer.Dial,
		TLSHandshakeTimeout: 16 * time.Second,
		DisableCompression:  true,
		DisableKeepAlives:   true,
	},
}

func getClient(skip string) *http.Client {
	/*_skip := false
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
	}*/
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
		next()
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
	req.Close = true
	resp, errResp := client.Do(req)
	if errResp != nil {
		next()
		return errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		next()
		return errResponse
	}

	for k := range resp.Header {
		newData.SetHeader(k, resp.Header.Get(k))
	}
	newData.SetHeader("status", fmt.Sprintf("%d", resp.StatusCode))
	newData.SetBody(string(data))
	out <- newData
	next()
	return nil
}
