package gonnie

//Header represent a key-value header pattern
type Header map[string]string

//Add new entry to map
func (h Header) Add(key, value string) {
	h[key] = value
}

//Get entry from map
func (h Header) Get(key string) string {
	return h[key]
}

//Del entry from map
func (h Header) Del(key string) {
	delete(h, key)
}

//ListKeys list keys from Map
func (h Header) ListKeys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}
