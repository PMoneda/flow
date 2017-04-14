package gonnie

type Processor func(Exchange) Exchange

// IRoute implements route builder pattern
type IRoute interface {
	From(string) IRoute
	To(string) IRoute
	SetHeader(string, interface{}) IRoute
	Processor(Processor) IRoute
	SetBody(interface{}) IRoute
	Log(string) IRoute
}

// Route is the struct that execute flow
type Route struct {
	context *Context
}

// From start point
func (r *Route) From(uri string) IRoute {
	return r
}

// To Next Point
func (r *Route) To(uri string) IRoute {
	return r
}

// Log message
func (r *Route) Log(msg string) IRoute {
	return r
}

// SetHeader in message
func (r *Route) SetHeader(key string, value interface{}) IRoute {
	return r
}

// SetBody in message
func (r *Route) SetBody(body interface{}) IRoute {
	return r
}

// Processor execute a function to process message
func (r *Route) Processor(p Processor) IRoute {
	return r
}

// NewRoute return a new and empty route
func NewRoute(c *Context) IRoute {
	r := Route{}
	r.context = c
	return &r
}
