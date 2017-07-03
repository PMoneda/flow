package flow

//HeaderMap represent a key-value header pattern
type HeaderMap map[string]string

//Add new entry to map
func (h HeaderMap) Add(key, value string) {
	h[key] = value
}

//Get entry from map
func (h HeaderMap) Get(key string) string {
	s, exist := h[key]
	if !exist {
		s = ""
	}
	return s
}

//Del entry from map
func (h HeaderMap) Del(key string) {
	delete(h, key)
}

//ListKeys list keys from Map
func (h HeaderMap) ListKeys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

//HeaderFnc is a type to execute logical conditions from head values
type HeaderFnc func(*Flow) (interface{}, *Flow)

//Header is the main entry point to execute logical conditions with header values
func Header(s string) HeaderFnc {
	return (func(pipe *Flow) (interface{}, *Flow) {
		header := pipe.GetHeader()
		return header.Get(s), pipe
	})
}

//IsEqualTo returns true if header[key] == s
func (h HeaderFnc) IsEqualTo(s string) HeaderFnc {
	return func(pipe *Flow) (interface{}, *Flow) {
		obj, _ := h(pipe)
		return obj.(string) == s, pipe
	}
}

//Exist returns true if header[key] exist
func (h HeaderFnc) Exist() HeaderFnc {
	return func(pipe *Flow) (interface{}, *Flow) {
		obj, _ := h(pipe)
		return obj.(string) != "", pipe
	}
}
