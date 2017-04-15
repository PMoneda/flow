package gonnie

import "testing"
import "strconv"
import "fmt"

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
