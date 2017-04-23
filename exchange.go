package gonnie

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

//Header represent a key-value header pattern
type Header map[string]string

//Add new entry to map
func (h Header) Add(key, value string) {
	h[key] = value
}

//Get entry from map
func (h Header) Get(key string) string {
	return h[key]
}

//ListKeys list keys from Map
func (h Header) ListKeys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

//Exchange is a middleware message used between process
type Exchange struct {
	inHead  Header
	in      *bytes.Buffer
	outHead Header
	out     *bytes.Buffer
}

// GetIn returns input message
func (e *Exchange) GetIn() *bytes.Buffer {
	return e.in
}

// GetOut returns output buffer
func (e *Exchange) GetOut() *bytes.Buffer {
	return e.out
}

// GetInHeader get input header
func (e *Exchange) GetInHeader() Header {
	return e.inHead
}

// GetOutHeader get output header
func (e *Exchange) GetOutHeader() Header {
	return e.outHead
}

//BindJSON binds json body to interface
func (e *Exchange) BindJSON(v interface{}) error {
	return json.Unmarshal(e.GetIn().Bytes(), &v)
}

//WriteJSON marshall interface to JSON and set body
func (e *Exchange) WriteJSON(v interface{}) error {
	x, err := json.Marshal(v)
	if err != nil {
		return err
	}
	e.GetOut().Write(x)
	return nil
}

//WriteXML marshall interface to XML and set body
func (e *Exchange) WriteXML(v interface{}) error {
	x, errXML := xml.Marshal(v)
	if errXML != nil {
		return errXML
	}
	e.GetOut().Write(x)
	return nil
}

//BindXML binds xml body to interface
func (e *Exchange) BindXML(v interface{}) error {
	return xml.Unmarshal(e.GetIn().Bytes(), &v)
}

// NewExchange creates new exchange message
func NewExchange() *Exchange {
	e := Exchange{
		inHead:  make(Header),
		outHead: make(Header),
		in:      bytes.NewBuffer(nil),
		out:     bytes.NewBuffer(nil),
	}
	return &e
}
