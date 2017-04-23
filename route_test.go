package gonnie

import "testing"
import "strconv"
import "fmt"

func TestShouldExecuteProcess(t *testing.T) {
	r := NewRouteWithContext()
	str := r.Processor(_sum).Processor(_sum).Processor(_sum).Body()
	if str != "3" {
		t.Fail()
	}
}
func _sum(ext *Exchange) {
	body := ext.GetIn().String()
	if body != "" {
		n, _ := strconv.Atoi(body)
		ext.GetOut().WriteString(fmt.Sprintf("%d", n+1))
	} else {
		ext.GetOut().WriteString("1")
	}

}

func TestShouldConsumeRestService(t *testing.T) {
	r := NewRouteWithContext()
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
	r = r.Processor(func(e *Exchange) {
		d := Root{}
		if e.BindJSON(&d) != nil || e.WriteXML(d) != nil {
			t.Fail()
		}
	})
	r = nil
}

func TestShouldConvertJSONInputMessageToOutPutUsingTransform(t *testing.T) {
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

	r := NewRouteWithContext()
	b := r.From("http://services.groupkt.com/state/get/IND/UP").Transform(from, "json", to).Body()
	if b != expected {
		t.Fail()
	}
}
