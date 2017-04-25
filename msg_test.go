package gonnie

import (
	"fmt"
	"testing"
)

const from = `
Header:
Content-Type: application/x-www-form-urlencoded
Cache-Control: no-cache

Body: 
grant_type=client_credentials&scope=cobranca.registro-boletos

`

func TestMsgUri(t *testing.T) {
	p := NewPipe()
	b := p.From("msg://", from).Body()
	fmt.Println(b.(string))
}
