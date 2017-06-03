package flow

//Choice is the type for conditional structure of a Pipe
type Choice struct {
	pipe     IPipe
	execute  bool
	executed bool
}

//NewChoice creates a new choice structure with a message
func NewChoice(p IPipe) *Choice {
	c := Choice{
		pipe:     p,
		execute:  false,
		executed: false,
	}
	return &c
}

//To execute a connector
func (c *Choice) To(url string, params ...interface{}) *Choice {
	if len(c.pipe.GetFails()) > 0 {
		return c
	}
	if c.execute {
		c.pipe = c.pipe.To(url, params...)
		c.executed = true
	}
	return c
}

// When is the conditional tester
func (c *Choice) When(e HeaderFnc) *Choice {
	if len(c.pipe.GetFails()) > 0 {
		return c
	}
	obj, _ := e(c.pipe)
	c.execute = false
	c.execute = obj.(bool)
	return c
}

//Otherwise is executed when others condiotions is not true - Like Else
func (c *Choice) Otherwise() *Choice {
	if len(c.pipe.GetFails()) > 0 {
		return c
	}
	c.execute = false
	if !c.executed {
		c.execute = true
	}
	return c
}
