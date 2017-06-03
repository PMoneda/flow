package flow

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

var letterFormat Transform = `Hello! How old are you? {{.age}}? \n-No, I'm {{.grade}}`

func TestShouldTransformDataXMLToLetter(t *testing.T) {
	transformed := data.TransformFromXML(inputFormat, letterFormat)
	if transformed != `Hello! How old are you? 10? \n-No, I'm 3` {
		t.Fail()
	}
}

func TestShouldTransformDataJSONToXML(t *testing.T) {
	var json Transform = `{ "first": 10, "second": 2, "third":{"a":1,"b":6,"c":7}, "bla":[1,2,3] }`
	var from Transform = `{ "first": "{{age}}", "second": "{{grade}}", "third":{"a":"{{test}}"} }`
	transformed := json.TransformFromJSON(from, letterFormat)
	if transformed != `Hello! How old are you? 10? \n-No, I'm 2` {
		t.Fail()
	}
}
