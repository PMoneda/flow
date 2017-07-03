package flow

import (
	"encoding/json"
	"encoding/xml"
	"errors"
)

//Transform date from template to other
type Transform string

//PipeProcessor is a function type to execute pipe workflow
type PipeProcessor func(*ExchangeMessage, Message, func())

//Message is a channel to change message between pipes
type Message chan interface{}

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
	case error:
		return t
	default:
		return errors.New("Invalid datatype")
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
	case error:
		return t
	default:
		return errors.New("Invalid datatype")
	}
}
