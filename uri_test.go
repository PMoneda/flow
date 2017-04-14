package gonnie

import "testing"
import "net/url"
import "fmt"

func TestShouldParseURL(t *testing.T) {
	u := "file://desktop.documents?name=teste.txt"
	url, err := url.Parse(u)
	if err != nil {
		t.Fail()
	}
	fmt.Println(url.Host)
	fmt.Println(url.RawQuery)
	fmt.Println(url.Scheme)

}
