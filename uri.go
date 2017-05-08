package gonnie

import "net/url"

import "sync"

type Uri struct {
	protocol string
	host     string
	options  url.Values
	raw      string
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
	"http":   httpConector,
	"https":  httpConector,
	"direct": directConector,
	"msg":    msg,
	"print":  printConector,
}

//RegisterConector register a new conector to use as From("my-connector://...")
func RegisterConector(name string, callback func(func(), *ExchangeMessage, Message, Uri, ...interface{}) error) {
	_lockConectors.Lock()
	defer _lockConectors.Unlock()
	pipeConectors[name] = callback
}
