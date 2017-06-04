package main

import (
	"fmt"

	"github.com/PMoneda/flow"
)

func main() {
	cityService()
}

func cityService() {
	fromOrignalAPI := `
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
	responseToMyAPI := `
		<country name="{{.country}}">
			<name>{{.name}}</name>
			<biggest_city>{{.biggest_city}}</biggest_city>
			<capital>{{.capital}}</capital>
		</country>
	`
	pipe := flow.NewPipe().From("http://services.groupkt.com/state/get/IND/UP")
	pipe.To("transform://?format=json", fromOrignalAPI, responseToMyAPI)
	fmt.Println(pipe.GetBody())
}
