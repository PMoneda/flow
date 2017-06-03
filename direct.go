package flow

func directConnector(next func(), e *ExchangeMessage, out Message, u URI, params ...interface{}) error {
	if len(params) > 0 {
		e.SetBody(params[0])
	}
	out <- e
	next()
	return nil
}
