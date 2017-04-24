package gonnie

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	Transform(Transform, string, Transform) IPipe
	Flush()
}

type Pipe struct {
	pipes Stack
}

func NewPipe() IPipe {
	p := Pipe{
		pipes: make(Stack, 0, 0),
	}
	return &p
}
func (p *Pipe) Body() interface{} {
	in := p.pipes.Top().(Message)
	for msg := range in {
		m := msg.(*ExchangeMessage)
		return m.body
	}
	return nil
}

//Processor on message
func (p *Pipe) Processor(proc PipeProcessor) IPipe {
	out := make(Message)
	in := p.pipes.Top().(Message)
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
	out := make(Message)
	in := p.pipes.Top().(Message)
	p.pipes.Push(out)
	go func() {
		for msg := range in {
			m := msg.(*ExchangeMessage)
			m.SetHead(k, v)
			out <- m
		}
		close(out)
	}()
	return p
}

//SetBody on message
func (p *Pipe) SetBody(b interface{}) IPipe {
	out := make(Message)
	in := p.pipes.Top().(Message)
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

func (p *Pipe) From(url string, params ...interface{}) IPipe {
	out := make(Message)
	p.pipes.Push(out)
	go func() {
		u, err := processURI(url)
		if err != nil {
			panic(err)
		}
		pipeConectors[u.protocol](func() {
			close(out)
		}, NewExchangeMessage(), out, u, params...)
	}()

	return p
}

func (p *Pipe) Transform(from Transform, mode string, to Transform) IPipe {
	out := make(Message)
	in := p.pipes.Top().(Message)
	p.pipes.Push(out)
	go func() {
		for msg := range in {
			m := msg.(*ExchangeMessage)
			t := Transform(m.body.(string))
			var trans string
			if "json" == mode {
				trans = string(t.TransformFromJSON(from, to))
			} else {
				trans = string(t.TransformFromXML(from, to))
			}
			m.body = trans
			out <- m
		}
		close(out)
	}()
	return p
}

func (p *Pipe) To(url string, params ...interface{}) IPipe {
	out := make(Message)
	in := p.pipes.Top().(Message)
	p.pipes.Push(out)
	go func() {
		u, _ := processURI(url)
		for n := range in {
			msg := n.(*ExchangeMessage)
			pipeConectors[u.protocol](func() {
				close(out)
			}, msg, out, u, params...)
		}

	}()
	return p
}

func (p *Pipe) Flush() {
	for v := range p.pipes.Top().(Message) {
		fmt.Println(v)
	}
	p.pipes.Clear()
}

type ExchangeMessage struct {
	head Header
	body interface{}
}

func (e *ExchangeMessage) SetHead(k, v string) {
	e.head.Add(k, v)
}
func (e *ExchangeMessage) GetHead(k string) string {
	return e.head.Get(k)
}

func (e *ExchangeMessage) DelHead(k string) {
	e.head.Del(k)
}

func (e *ExchangeMessage) SetBody(b interface{}) {
	e.body = b
}
func (e *ExchangeMessage) GetBody() interface{} {
	return e.body
}

func NewExchangeMessage() *ExchangeMessage {
	e := ExchangeMessage{
		head: make(Header),
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
