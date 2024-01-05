package templates

import (
	"bytes"
	"embed"
	"text/template"
)

const mainTemplate = "events.gotmpl"

var (
	//go:embed *.gotmpl
	templates     embed.FS
	templateFuncs = template.FuncMap{
		"isError": isError,
	}
	templatesParsed = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templates, "*"))
)

func Exec(d *Dot) ([]byte, error) {
	var buf bytes.Buffer
	if err := templatesParsed.ExecuteTemplate(&buf, mainTemplate, d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func isError(ev []*Event, isErr bool) any {
	return &struct {
		IsError bool
		Events  []*Event
	}{
		IsError: isErr,
		Events:  ev,
	}
}

type (
	Dot struct {
		Package  string
		Imports  []string
		Events   []*Event
		Verbosef []*Event
		Errorf   []*Event
		Custom   []*Event
	}
	Event struct {
		Name   string
		Type   string
		Level  string
		Nice   string
		Format string
		Args   []Arg
		Custom bool
	}
	Arg struct {
		Name string
		Type string
	}
)
