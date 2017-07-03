package flow

import "errors"

func unmarshallConnector(e *ExchangeMessage, u URI, params ...interface{}) error {
	if len(params) < 1 {
		err := errors.New("You should give a objeto to unmarshall to")
		e.body = err
		e.SetHeader("error", err.Error())
		return err
	}
	if u.GetOption("format") == "json" {
		err := e.BindJSON(params[0])
		if err != nil {
			e.body = err
			e.SetHeader("error", err.Error())
			return err
		}
		e.SetBody(params[0])
	} else if u.GetOption("format") == "xml" {
		err := e.BindXML(params[0])
		if err != nil {
			e.body = err
			e.SetHeader("error", err.Error())
			return err
		}
		e.SetBody(params[0])
	} else {
		err := errors.New("Invalid format type")
		e.body = err
		e.SetHeader("error", err.Error())
		return err
	}
	return nil
}
