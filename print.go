package flow

import (
	"fmt"
)

func printConnector(e *ExchangeMessage, u URI, params ...interface{}) error {
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
	return nil
}
