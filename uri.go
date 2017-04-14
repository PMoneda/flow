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

func exec(ctx *Context, u uri) {

}
