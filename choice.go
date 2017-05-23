package gonnie

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

func (c *Choice) When(e HeaderFnc) *Choice {
	if len(c.pipe.GetFails()) > 0 {
		return c
	}
	obj, _ := e(c.pipe)
	c.execute = false
	c.execute = obj.(bool)
	return c
}

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
