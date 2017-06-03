package flow

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
)

//Transform date from template to other
type Transform string

//PipeProcessor is a function type to execute pipe workflow
type PipeProcessor func(*ExchangeMessage, Message, func())

//Message is a channel to change message between pipes
type Message chan interface{}

//IPipe is the interface with Pipe support
type IPipe interface {
	From(string, ...interface{}) IPipe
	To(string, ...interface{}) IPipe
	SetHeader(string, string) IPipe
	SetBody(interface{}) IPipe
	Processor(PipeProcessor) IPipe
	Body() interface{}
	GetBody() interface{}
	GetFails() []error
	GetHeader() HeaderMap
	Header() HeaderMap
	Transform(Transform, string, Transform, template.FuncMap) IPipe
	Choice() *Choice
}

// Pipe is the main data structure of Flow Pipe controls all the execution flow
type Pipe struct {
	pipes Stack
	fails []error
}

// NewPipe creates a new empty Pipe
func NewPipe() IPipe {
	p := Pipe{
		pipes: make(Stack, 0, 0),
		fails: make([]error, 0, 0),
	}
	return &p
}

//Choice is the conditional flow of a pipe
func (p *Pipe) Choice() *Choice {
	if len(p.fails) > 0 {
		printFails(p)
		return NewChoice(p)
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	input := <-in
	go func() {
		out <- input
	}()
	return NewChoice(p)
}

//GetBody get final body result from executed flow
func (p *Pipe) GetBody() interface{} {
	in := p.pipes.Pop().(Message)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	if len(p.fails) > 0 {
		printFails(p)
		return p.fails
	}

	for msg := range in {
		m := msg.(*ExchangeMessage)
		return m.body
	}
	return nil
}

//GetHeader same as GetBody but instead of return body of message this method returns the message header
func (p *Pipe) GetHeader() HeaderMap {
	in := p.pipes.Pop().(Message)
	for msg := range in {
		m := msg.(*ExchangeMessage)
		return m.head
	}
	return nil
}

//Body returns body but doesn't lock the flow
func (p *Pipe) Body() interface{} {
	if len(p.fails) > 0 {
		printFails(p)
		return p.fails
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	var body interface{}
	for msg := range in {
		m := msg.(*ExchangeMessage)
		body = m.body
		go func() {
			out <- m
			close(out)
		}()

	}
	return body
}

//Header return current Header but doesn't lock the flow
func (p *Pipe) Header() HeaderMap {

	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	var h HeaderMap
	for msg := range in {
		m := msg.(*ExchangeMessage)
		h = m.head
		go (func(ms *ExchangeMessage) {
			out <- ms
			close(out)
		})(m)
		break
	}
	return h
}

//Processor execute a user-pass go function in the flow
func (p *Pipe) Processor(proc PipeProcessor) IPipe {

	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	go func() {
		for msg := range in {
			m := msg.(*ExchangeMessage)
			proc(m, out, func() {
				close(out)
			})
		}
	}()
	return p
}

//SetHeader on message
func (p *Pipe) SetHeader(k, v string) IPipe {

	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	go func() {
		for msg := range in {
			m := msg.(*ExchangeMessage)
			m.SetHeader(k, v)
			out <- m
		}
		close(out)
	}()
	return p
}

//SetBody on message
func (p *Pipe) SetBody(b interface{}) IPipe {

	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	go func() {

		for msg := range in {
			m := msg.(*ExchangeMessage)
			m.SetBody(b)
			out <- m
		}
		close(out)
	}()
	return p
}

//From is the entrypoint of a flow, all flows need to be started with From
func (p *Pipe) From(url string, params ...interface{}) IPipe {
	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	p.pipes.Push(out)
	go func() {
		u, err := processURI(url)
		if err != nil {
			close(out)
			return
		}
		pipeConnectors[u.protocol](func() {
			close(out)
		}, NewExchangeMessage(), out, u, params...)
	}()

	return p
}

//Transform -deprecated
func (p *Pipe) Transform(from Transform, mode string, to Transform, fncs template.FuncMap) IPipe {
	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	go func() {
		for msg := range in {
			m := msg.(*ExchangeMessage)
			t := Transform(m.body.(string))
			var trans string
			var s string
			var err error
			if "json" == mode {
				s, err = t.TransformFromJSON(from, to, fncs)
			} else {
				s, err = t.TransformFromXML(from, to, fncs)
			}
			if err != nil {
				p.fails = append(p.fails, err)
			}
			trans = string(s)
			m.body = trans
			out <- m
		}
		close(out)
	}()
	return p
}

//To is a method to create a flow based on connectors
func (p *Pipe) To(url string, params ...interface{}) IPipe {
	if len(p.fails) > 0 {
		printFails(p)
		return p
	}
	out := make(Message)
	in := p.pipes.Pop().(Message)
	p.pipes.Push(out)
	go func() {
		u, errURI := processURI(url)
		if errURI != nil {
			p.fails = append(p.fails, errURI)
			close(out)
			return
		}
		for n := range in {
			msg := n.(*ExchangeMessage)
			err := pipeConnectors[u.protocol](func() {}, msg, out, u, params...)
			if err != nil {
				fmt.Printf("Erro: %s\nURI:%s\n", err, url)
				p.fails = append(p.fails, err)
				close(out)
				return
			}
			return
		}
	}()

	return p
}
func printFails(p *Pipe) {
	for _, e := range p.fails {
		fmt.Println(e)
	}
}

//GetFails returns all error occured dURIng execution flow
func (p *Pipe) GetFails() []error {
	return p.fails
}

//ExchangeMessage is the message exchange inner Pipe
type ExchangeMessage struct {
	head HeaderMap
	body interface{}
}

//SetHeader add a key-value header to a message
func (e *ExchangeMessage) SetHeader(k, v string) {
	e.head.Add(k, v)
}

//GetHeader returns a header value based on key
func (e *ExchangeMessage) GetHeader(k string) string {
	return e.head.Get(k)
}

//DelHeader delete header from message
func (e *ExchangeMessage) DelHeader(k string) {
	e.head.Del(k)
}

//GetHeaderMap return the map from  header
func (e *ExchangeMessage) GetHeaderMap() HeaderMap {
	return e.head
}

//ClearHeader removes all entries from header
func (e *ExchangeMessage) ClearHeader() {
	keys := e.head.ListKeys()
	for _, k := range keys {
		e.head.Del(k)
	}
}

//SetBody writes the body to message
func (e *ExchangeMessage) SetBody(b interface{}) {
	e.body = b
}

//GetBody return body from message
func (e *ExchangeMessage) GetBody() interface{} {
	return e.body
}

//NewExchangeMessage creates new empty ExchangeMessage
func NewExchangeMessage() *ExchangeMessage {
	e := ExchangeMessage{
		head: make(HeaderMap),
	}
	return &e
}

//BindJSON binds json body to interface
func (e *ExchangeMessage) BindJSON(v interface{}) error {
	switch t := e.body.(type) {
	case string:
		return json.Unmarshal([]byte(t), &v)
	case []byte:
		return json.Unmarshal(t, &v)
	default:
		panic("Invalid datatype")
	}
}

//WriteJSON marshall interface to JSON and set body
func (e *ExchangeMessage) WriteJSON(v interface{}) error {
	x, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.body = x
	return nil
}

//WriteXML marshall interface to XML and set body
func (e *ExchangeMessage) WriteXML(v interface{}) error {
	x, errXML := xml.Marshal(v)
	if errXML != nil {
		return errXML
	}
	e.body = x
	return nil
}

//BindXML binds xml body to interface
func (e *ExchangeMessage) BindXML(v interface{}) error {
	switch t := e.body.(type) {
	case string:
		return xml.Unmarshal([]byte(t), &v)
	case []byte:
		return xml.Unmarshal(t, &v)
	default:
		panic("Invalid datatype")
	}
}
