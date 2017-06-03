package flow

import "errors"

func unmarshallConnector(next func(), e *ExchangeMessage, out Message, u URI, params ...interface{}) error {
	if len(params) < 1 {
		return errors.New("You should give a objeto to unmarshall to")
	}
	if u.GetOption("format") == "json" {
		err := e.BindJSON(params[0])
		if err != nil {
			return err
		}
		e.SetBody(params[0])
	} else if u.GetOption("format") == "xml" {
		err := e.BindXML(params[0])
		if err != nil {
			return err
		}
		e.SetBody(params[0])
	} else {
		return errors.New("Invalid format type")
	}

	out <- e
	next()
	return nil
}
