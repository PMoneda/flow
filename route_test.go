package gonnie

import "testing"
import "strconv"
import "fmt"

func TestShouldCreateAWorkFlow(t *testing.T) {
	c := NewContext()
	r := NewRoute(c)
	r = r.From("direct://hello").To("http://somesite.com").To("file://destiny?name=ggg.txt")
}

func TestShouldExecuteProcess(t *testing.T) {
	c := NewContext()
	r := NewRoute(c)
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
