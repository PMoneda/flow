package flow

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var cfg = &tls.Config{
	InsecureSkipVerify: true,
}
var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:     cfg,
		Dial:                defaultDialer.Dial,
		TLSHandshakeTimeout: 16 * time.Second,
		DisableCompression:  true,
		DisableKeepAlives:   true,
	},
}

func getClient(skip, timeout string) *http.Client {
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

	t, _ := strconv.Atoi(timeout)
	tOut := time.Duration(t)
	client.Timeout = tOut * time.Second

	return client
}
func httpConnector(e *ExchangeMessage, u URI, params ...interface{}) error {
	var skip string
	var opts map[string]string
	method := "GET"
	authMethod := ""
	username := ""
	password := ""
	timeout := ""
	if len(params) > 0 {
		opts = params[0].(map[string]string)
		skip = opts["insecureSkipVerify"]
		timeout = opts["timeout"]
		method = opts["method"]
		authMethod = opts["auth"]
		username = opts["username"]
		password = opts["password"]
	}
	client := getClient(skip, timeout)
	var req *http.Request
	var err error
	b := e.GetBody()
	var body io.Reader
	switch t := b.(type) {
	case string:
		body = strings.NewReader(t)
	default:
		if j, err := json.Marshal(b); err != nil {
			e.SetHeader("error", err.Error())
			e.SetBody(err)
		} else {
			body = strings.NewReader(string(j))
		}
	}
	req, err = http.NewRequest(method, u.raw, body)

	if err != nil {
		e.SetHeader("error", err.Error())
		e.SetBody(err)
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
		e.SetHeader("error", errResp.Error())
		e.SetBody(errResp)
		return errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		e.SetHeader("error", errResponse.Error())
		e.SetBody(errResponse)
		return errResponse
	}

	for k := range resp.Header {
		e.SetHeader(k, resp.Header.Get(k))
	}
	e.SetHeader("status", fmt.Sprintf("%d", resp.StatusCode))
	e.SetBody(string(data))
	return nil
}
