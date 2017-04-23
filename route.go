package gonnie

// Processor is the template for process callback
type Processor func(*Exchange)

//Transform date from template to other
type Transform string

// IRoute implements route builder pattern
type IRoute interface {
	From(string) IRoute
	To(string) IRoute
	SetHeader(string, string) IRoute
	Processor(Processor) IRoute
	SetBody(string) IRoute
	Log(string) IRoute
	Body() string
	Header() Header
	Transform(Transform, string, Transform) IRoute
}

// Route is the struct that execute flow
type Route struct {
	context *Context
}

// From start point
func (r *Route) From(uri string) IRoute {
	return r.execSimple(uri)
}

// To Next Point
func (r *Route) To(uri string) IRoute {
	return r.execSimple(uri)
}

// Log message
func (r *Route) Log(msg string) IRoute {
	log(r.context, msg)
	return r
}

// SetHeader in message
func (r *Route) SetHeader(key string, value string) IRoute {
	r.context.GetMessage().GetInHeader().Add(key, value)
	return r
}

// SetBody in message
func (r *Route) SetBody(body string) IRoute {
	r.context.GetMessage().GetIn().WriteString(body)
	return r
}

func (r *Route) execSimple(p ...string) IRoute {
	err := execURI(r.context, p...)
	if err != nil {
		r.context.stack.Pop()
		panic(err)
	}
	return r
}

// Processor execute a function to process message
func (r *Route) Processor(p Processor) IRoute {
	r.context.PushMessage()
	p(r.context.GetMessage())
	return r
}

// Body get body from message
func (r *Route) Body() string {
	return r.context.GetMessage().GetOut().String()
}

// Header get header from message
func (r *Route) Header() Header {
	return r.context.GetMessage().GetOutHeader()
}

//Transform data from template A to template B
func (r *Route) Transform(from Transform, mode string, to Transform) IRoute {
	r.context.PushMessage()
	t := Transform(r.context.GetMessage().GetIn().String())
	var trans string
	if "json" == mode {
		trans = string(t.TransformFromJSON(from, to))
	} else {
		trans = string(t.TransformFromXML(from, to))
	}
	r.context.GetMessage().GetOut().WriteString(trans)
	return r
}

// NewRoute return a new and empty route
func NewRoute(c *Context) IRoute {
	r := Route{}
	r.context = c
	return &r
}
