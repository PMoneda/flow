package flow

func setConnector(e *ExchangeMessage, u URI, params ...interface{}) error {
	if len(params) > 0 {
		if u.options.Get("prop") == "body" {
			e.SetBody(params[0])
		} else if u.options.Get("prop") == "header" {
			mp := params[0].(map[string]string)
			for k, v := range mp {
				e.head.Add(k, v)
			}
		}
	}
	return nil
}
