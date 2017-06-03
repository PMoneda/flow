package flow

import (
	"net/url"
	"sync"
)

type Uri struct {
	protocol string
	host     string
	options  url.Values
	raw      string
}

//GetOption from URI
func (u Uri) GetOption(k string) string {
	return u.options.Get(k)
}

func (u Uri) GetProtocol() string {
	return u.protocol
}

func (u Uri) GetHost() string {
	return u.host
}

func (u Uri) GetRaw() string {
	return u.raw
}

func processURI(u string) (Uri, error) {
	url, err := url.Parse(u)
	if err != nil {
		return Uri{}, err
	}
	ur := Uri{}
	ur.host = url.Host
	ur.protocol = url.Scheme
	ur.options = url.Query()
	ur.raw = u
	return ur, nil
}

var _lockConectors sync.Mutex
var pipeConectors = map[string]func(func(), *ExchangeMessage, Message, Uri, ...interface{}) error{
	"http":       httpConector,
	"https":      httpConector,
	"direct":     directConector,
	"print":      printConector,
	"transform":  transformConector,
	"template":   templateConector,
	"set":        setConector,
	"message":    messageConector,
	"unmarshall": unmarshallConector,
}

//RegisterConector register a new conector to use as From("my-connector://...")
func RegisterConector(name string, callback func(func(), *ExchangeMessage, Message, Uri, ...interface{}) error) {
	_lockConectors.Lock()
	defer _lockConectors.Unlock()
	pipeConectors[name] = callback
}
