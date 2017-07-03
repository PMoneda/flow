package flow

func directConnector(e *ExchangeMessage, u URI, params ...interface{}) error {
	if len(params) > 0 {
		e.SetBody(params[0])
	}
	return nil
}
