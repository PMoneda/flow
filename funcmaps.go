package gonnie

import (
	"html/template"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"today":   today,
	"brdate":  brDate,
	"replace": replace,
	"trim":    trim,
}

func today() time.Time {
	return time.Now()
}

func brDate(d time.Time) string {
	return d.Format("02/01/2006")
}

func replace(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
