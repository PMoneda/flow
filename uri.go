package flow

import (
	"net/url"
	"sync"
)

//URI is the pattern of connectors
type URI struct {
	protocol string
	host     string
	options  url.Values
	raw      string
}

//GetOption from URI
func (u URI) GetOption(k string) string {
	return u.options.Get(k)
}

//GetProtocol return connector schema
func (u URI) GetProtocol() string {
	return u.protocol
}

//GetHost retruns connector host
func (u URI) GetHost() string {
	return u.host
}

//GetRaw returns a plain text connector url
func (u URI) GetRaw() string {
	return u.raw
}

func processURI(u string) (URI, error) {
	url, err := url.Parse(u)
	if err != nil {
		return URI{}, err
	}
	ur := URI{}
	ur.host = url.Host
	ur.protocol = url.Scheme
	ur.options = url.Query()
	ur.raw = u
	return ur, nil
}

var _lockConectors sync.Mutex
var pipeConectors = map[string]func(func(), *ExchangeMessage, Message, URI, ...interface{}) error{
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
func RegisterConector(name string, callback func(func(), *ExchangeMessage, Message, URI, ...interface{}) error) {
	_lockConectors.Lock()
	defer _lockConectors.Unlock()
	pipeConectors[name] = callback
}
