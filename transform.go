package gonnie

import (
	"bytes"
	"html/template"
	"strings"

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
