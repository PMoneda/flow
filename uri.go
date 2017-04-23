package gonnie

import "net/url"
import "fmt"
import "sync"

type uri struct {
	protocol string
	host     string
	options  url.Values
	raw      string
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
	ur.raw = u
	return ur, nil
}

var _lockConectors sync.Mutex

var execMux = map[string]func(*Context, uri, ...string) error{
	"direct": direct,
	"http":   _http,
	"https":  _https,
	"file":   file,
}

func execURI(ctx *Context, u ...string) error {
	ur, err := processURI(u[0])
	if err != nil {
		return err
	}
	errExec := execMux[ur.protocol](ctx, ur, u...)
	if errExec != nil {
		fmt.Println(u)
		fmt.Println(errExec.Error())
	}
	return nil
}

//RegisterConector register a new conector to use as From("my-connector://...")
func RegisterConector(name string, callback func(ctx *Context, u uri, s ...string) error) {
	_lockConectors.Lock()
	defer _lockConectors.Unlock()
	execMux[name] = callback
}
