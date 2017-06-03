package flow

import (
	"testing"
	"time"
)

func TestPipeProcessor(t *testing.T) {
	p := NewPipe()
	final := p.From("direct://a").Processor(func(m *ExchangeMessage, out Message, next func()) {
		go func() {
			time.Sleep(1 * time.Second)
			m.SetBody(1)
			out <- m
			next()
		}()
	}).Processor(func(m *ExchangeMessage, out Message, next func()) {
		go func() {
			time.Sleep(1 * time.Second)
			m.SetBody(m.GetBody().(int) + 1)
			out <- m
			next()
		}()
	}).Processor(func(m *ExchangeMessage, out Message, next func()) {
		go func() {
			time.Sleep(1 * time.Second)
			m.SetBody(m.GetBody().(int) + 1)
			out <- m
			next()
		}()
	}).Body()

	if 3 != final.(int) {
		t.Fail()
	}
}

func TestShouldConsumeRestService(t *testing.T) {
	r := NewPipe()
	type Root struct {
		Response struct {
			Messages []string `json:"messages"`
			Result   []struct {
				Name        string `json:"name"`
				Alpha2_Code string `json:"alpha2_code"`
				Alpha3_Code string `json:"alpha3_code"`
			} `json:"result"`
		} `json:"RestResponse"`
	}
	r = r.From("http://services.groupkt.com/country/get/all")
	r = r.Processor(func(e *ExchangeMessage, out Message, next func()) {
		d := Root{}
		if e.BindJSON(&d) != nil || e.WriteXML(d) != nil {
			t.Fail()
		}
		out <- e
		next()
	})
	r = nil
}

func TestShouldConvertJSONInputMessageToOutPutUsingTransformOnPipe(t *testing.T) {
	var from Transform = `
		{
			"RestResponse" : {				
				"result" : {
					"country" : "{{country}}",
					"name" : "{{name}}",
					"largest_city" : "{{biggest_city}}",
					"capital" : "{{capital}}"
				}
			}
		}
	`
	var to Transform = `
		<country name="{{.country}}">
			<name>{{.name}}</name>
			<biggest_city>{{.biggest_city}}</biggest_city>
			<capital>{{.capital}}</capital>
		</country>
	`
	var expected = `
		<country name="IND">
			<name>Uttar Pradesh</name>
			<biggest_city>Kanpur</biggest_city>
			<capital>Lucknow</capital>
		</country>
	`

	r := NewPipe()
	b := r.From("http://services.groupkt.com/state/get/IND/UP").Transform(from, "json", to, nil).Body()
	if b != expected {
		t.Fail()
	}
}
