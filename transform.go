package gonnie

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/beevik/etree"
)

//TransformFromXML data from A to B
func (original Transform) TransformFromXML(from, to Transform) string {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(string(original)); err != nil {
		panic(err)
	}
	translator := getTranslator(string(from))
	values := createMapValues(translator, doc)
	return generateTransformedMessage(values, to)
}

//TransformFromJSON data from A to B
func (original Transform) TransformFromJSON(from, to Transform) string {
	jsonParsed, _ := gabs.ParseJSON([]byte(from))
	ch, _ := jsonParsed.ChildrenMap()
	path := ""
	values := make(map[string]string)
	walkTreeJSON(ch, &path, values)
	oriParsed, _ := gabs.ParseJSON([]byte(original))
	parsedMap := make(map[string]string)
	for k, v := range values {
		s1 := strings.Replace(v, "\"{{", "", -1)
		s1 = strings.Replace(s1, "}}\"", "", -1)
		parsedMap[s1] = oriParsed.Path(k).String()

	}
	return generateTransformedMessage(parsedMap, to)
}

func generateTransformedMessage(values map[string]string, to Transform) string {
	buf := bytes.NewBuffer(nil)
	t := template.Must(template.New("transform").Funcs(funcMap).Parse(string(to)))
	err := t.ExecuteTemplate(buf, "transform", values)
	if err != nil {
		panic(err)
	}
	return buf.String()
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
	}
}

func walkTreeJSON(ch map[string]*gabs.Container, path *string, values map[string]string) {

	for k, v := range ch {
		if *path == "" {
			*path = k
		} else {
			*path = *path + "." + k
		}

		chd, err := v.ChildrenMap()
		if err == nil {
			walkTreeJSON(chd, path, values)
		} else {
			values[*path] = v.String()
			*path = ""
		}

	}
}
