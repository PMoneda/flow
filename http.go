package gonnie

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strconv"
)

func _http(ctx *Context, u uri, s ...string) error {
	return _https(ctx, u, s...)
}

func _https(ctx *Context, u uri, s ...string) error {
	skip := u.options.Get("insecureSkipVerify")
	client := getClient(skip)
	req, err := http.NewRequest(u.options.Get("method"), u.options.Get("url"), ctx.GetMessage().GetIn())
	if err != nil {
		return err
	}
	authMethod := u.options.Get("auth")
	if authMethod == "basic" {
		req.SetBasicAuth(u.options.Get("user"), u.options.Get("password"))
	}
	header := ctx.GetMessage().GetInHeader()
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
	for k := range resp.Header {
		ctx.GetMessage().GetOutHeader().Add(k, resp.Header.Get(k))
	}
	_, errWrite := ctx.GetMessage().GetOut().Write(data)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

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
