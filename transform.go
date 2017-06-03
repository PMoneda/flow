package flow

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"html/template"
	"strings"
	"sync"

	"github.com/Jeffail/gabs"
	"github.com/beevik/etree"
)

func transformConector(next func(), m *ExchangeMessage, out Message, u Uri, params ...interface{}) error {
	t := Transform(m.body.(string))

	var trans string
	var s string
	var errFmt error
	var fncs template.FuncMap
	if len(params) > 2 && params[2] != nil {
		fncs = params[2].(template.FuncMap)
	}
	if "json" == u.options.Get("format") {
		s, errFmt = t.TransformFromJSON(convertToTransform(params[0]), convertToTransform(params[1]), fncs)
	} else {
		s, errFmt = t.TransformFromXML(convertToTransform(params[0]), convertToTransform(params[1]), fncs)
	}
	if errFmt != nil {
		return errFmt
	}
	trans = string(s)
	m.body = trans
	out <- m
	return nil
}

func convertToTransform(s interface{}) Transform {
	switch p := s.(type) {
	case string:
		return Transform(p)
	case Transform:
		return p
	}
	return ""
}

//TransformFromXML data from A to B
func (original Transform) TransformFromXML(from, to Transform, fncMap template.FuncMap) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(string(original)); err != nil {
		panic(err)
	}
	translator := getTranslator(string(from))
	values := createMapValues(translator, doc)
	return generateTransformedMessage(values, to, fncMap)
}

//TransformFromJSON data from A to B
func (original Transform) TransformFromJSON(from, to Transform, fncMap template.FuncMap) (string, error) {
	jsonParsed, _ := gabs.ParseJSON([]byte(from))
	ch, _ := jsonParsed.ChildrenMap()
	path := ""
	values := make(map[string]string)
	walkTreeJSON(ch, &path, values)
	oriParsed, _ := gabs.ParseJSON([]byte(original))
	parsedMap := make(map[string]string)
	for k, v := range values {
		parsedMap[cleanString(v)] = cleanString(oriParsed.Path(k).String())
	}
	return generateTransformedMessage(parsedMap, to, fncMap)
}
func cleanString(s string) string {
	s1 := strings.Replace(s, "{{", "", -1)
	s1 = strings.Replace(s1, "}}", "", -1)
	s1 = strings.Replace(s1, "\"", "", -1)
	s1 = strings.Replace(s1, "&#34;", "", -1)
	return s1
}

var tmpl = make(map[string]*template.Template)

func sha(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	sEnc := base64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))
	return sEnc
}

var ltmpl sync.Mutex

func generateTransformedMessage(values map[string]string, to Transform, fncMap template.FuncMap) (string, error) {
	buf := bytes.NewBuffer(nil)
	//var t *template.Template
	hash := sha(string(to))
	t, ok := tmpl[hash]
	if !ok {
		ltmpl.Lock()
		if fncMap != nil {
			t = template.Must(template.New(hash).Funcs(fncMap).Parse(string(to)))
		} else {
			t = template.Must(template.New(hash).Parse(string(to)))
		}
		tmpl[hash] = t
		ltmpl.Unlock()
	}
	err := t.ExecuteTemplate(buf, hash, values)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func createMapValues(translator map[string]string, doc *etree.Document) map[string]string {
	values := make(map[string]string)
	for k := range translator {
		for _, t := range doc.FindElements(k) {
			values[translator[k]] = t.Text()
		}
	}
	return values
}
func getTranslator(from string) map[string]string {
	docFrom := etree.NewDocument()
	if err := docFrom.ReadFromString(string(from)); err != nil {
		panic(err)
	}
	translator := make(map[string]string)
	for _, node := range docFrom.ChildElements() {
		walkTree(node, ""+node.Tag, translator)
	}
	return translator
}

func walkTree(elem *etree.Element, path string, translator map[string]string) {
	old := path
	if len(elem.ChildElements()) == 0 {
		if strings.HasPrefix(elem.Text(), "{{") {
			_translateTo := strings.Replace(elem.Text(), "{{", "", -1)
			_translateTo = strings.Replace(_translateTo, "}}", "", -1)
			translator[path] = _translateTo
		}
		return
	}
	for _, node := range elem.ChildElements() {
		path = path + "/" + node.Tag
		walkTree(node, path, translator)
		path = old

	}
}

func walkTreeJSON(ch map[string]*gabs.Container, path *string, values map[string]string) {
	for k, v := range ch {
		var old string
		if *path == "" {
			*path = k
		} else {
			old = *path
			*path = *path + "." + k
		}
		chd, err := v.ChildrenMap()
		if err == nil {
			walkTreeJSON(chd, path, values)
		} else {
			values[*path] = v.String()
		}
		*path = old

	}
}
