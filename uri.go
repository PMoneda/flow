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

var _lockConnectors sync.Mutex

var pipeConnectorsSync = map[string]func(*ExchangeMessage, URI, ...interface{}) error{
	"http":       httpConnector,
	"https":      httpConnector,
	"direct":     directConnector,
	"print":      printConnector,
	"transform":  transformConnector,
	"template":   templateConnector,
	"set":        setConnector,
	"message":    messageConnector,
	"unmarshall": unmarshallConnector,
	"crawler":    crawlerConnector,
}

//RegisterConnector register a new Connector to use as From("my-connector://...")
func RegisterConnector(name string, callback func(*ExchangeMessage, URI, ...interface{}) error) {
	_lockConnectors.Lock()
	defer _lockConnectors.Unlock()
	pipeConnectorsSync[name] = callback
}
