package flow

import (
	"fmt"
)

func printConnector(next func(), e *ExchangeMessage, out Message, u URI, params ...interface{}) error {
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
