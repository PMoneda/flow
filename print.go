package flow

import (
	"fmt"
)

func printConector(next func(), e *ExchangeMessage, out Message, u Uri, params ...interface{}) error {
	msg := u.options.Get("msg")
	//TODO refactor
	if msg == "${body}" {
		fmt.Println(e.body)
	} else if msg == "${head}" {
		for k, v := range e.head {
			fmt.Println(k + ":" + v)
		}
	} else {
		fmt.Println(msg)
	}

	out <- e
	next()
	return nil
}
