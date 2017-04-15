package gonnie

import "net/url"

type uri struct {
	protocol string
	host     string
	options  url.Values
}

func processURI(u string) (uri, error) {
	url, err := url.Parse(u)
	if err != nil {
		return uri{}, err
	}
	ur := uri{}
	ur.host = url.Host
	ur.protocol = url.Scheme
	ur.options = url.Query()
	return ur, nil
}

var execMux = map[string]func(*Context, ...string){
	"direct": direct,
	"http":   http,
	"file":   file,
}

func execURI(ctx *Context, u ...string) error {
	ur, err := processURI(u[0])
	if err != nil {
		return err
	}
	execMux[ur.protocol](ctx, u...)
	return nil
}
