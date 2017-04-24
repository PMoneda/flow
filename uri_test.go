package gonnie

import "testing"

func TestShouldParseURL(t *testing.T) {
	u := "file://desktop.documents?name=test.txt"
	ur, err := processURI(u)
	if err != nil {
		t.Fail()
	}

	if ur.protocol != "file" {
		t.Fail()
	}
}
