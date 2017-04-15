package gonnie

import "testing"

var data Transform = `<book><a>10</a></book><b>3</b>`
var inputFormat Transform = `<book><a>{{age}}</a></book><b>{{grade}}</b>`
var outPutFormat Transform = `{ "a":"{{.age}}","b":"{{.grade}}"}`

func TestShouldTransformData(t *testing.T) {
	transformed := data.TransformFromXML(inputFormat, outPutFormat)
	if transformed != `{ "a":"10","b":"3"}` {
		t.Fail()
	}
}
