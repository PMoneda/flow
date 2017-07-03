package flow

// Flow is the main data structure of Flow Pipe controls all the execution flow
type Flow struct {
	message *ExchangeMessage
	err     error
}

// NewFlow creates a new empty Flow
func NewFlow() *Flow {
	p := Flow{
		message: NewExchangeMessage(),
	}
	return &p
}

//From is the entrypoint of a flow, all flows need to be started with From
func (p *Flow) From(url string, params ...interface{}) *Flow {
	u, err := processURI(url)
	if err != nil {
		p.err = err
		return p
	}
	if err := pipeConnectorsSync[u.protocol](p.message, u, params...); err != nil {
		p.err = err
	}
	return p
}

//To is a method to create a flow based on connectors
func (p *Flow) To(url string, params ...interface{}) *Flow {
	if p.err != nil {
		return p
	}
	if u, errURI := processURI(url); errURI != nil {
		p.err = errURI
		return p
	} else if err := pipeConnectorsSync[u.protocol](p.message, u, params...); err != nil {
		p.err = err
	}
	return p
}

//SetBody to a execution Flow
func (p *Flow) SetBody(b interface{}) *Flow {
	if p.err != nil {
		return p
	}
	p.message.SetBody(b)
	return p
}

//SetHeader on message
func (p *Flow) SetHeader(k, v string) *Flow {
	if p.err != nil {
		return p
	}
	p.message.SetHeader(k, v)
	return p
}

//GetHeader from a message
func (p *Flow) GetHeader() HeaderMap {
	return p.message.head
}

//GetBody from a message
func (p *Flow) GetBody() interface{} {
	return p.message.body
}

//Choice builds a choice block
func (p *Flow) Choice() *Choice {
	return NewChoice(p)
}
