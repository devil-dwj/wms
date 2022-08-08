package main

import (
	"bytes"
	"text/template"
)

var tpl = `
{{ range .Errors }}

func Error{{ .CamelValue }}(format string, args ...interface{}) error {
	return errors.Errorf({{ .Status }}, int({{ .Name }}_{{ .CodeDsc }}), fmt.Sprintf(format, args...))
}

{{- end }}
`

type info struct {
	Name       string
	CamelValue string
	Status     int
	CodeDsc    string
}

type err struct {
	Errors []*info
}

func (e *err) Execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("Errors").Parse(tpl)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
