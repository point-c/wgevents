package templates

import (
	"embed"
	"fmt"
	helpers "github.com/point-c/generator-helpers"
	"strings"
	"text/template"
)

const (
	mainTemplate = "events.gotmpl"
	testTemplate = "events_test.gotmpl"
)

var (
	//go:embed *.gotmpl
	templates     embed.FS
	templateFuncs = template.FuncMap{
		"isError":     isError,
		"nilSlice":    nilSlice,
		"joinStr":     joinStr,
		"inc":         increment,
		"add":         add,
		"mul":         mul,
		"appendSlice": appendSlice,
		"consSlice":   consSlice,
	}
	templatesParsed = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templates, "*"))
)

func GetTemplate(test bool) string {
	if test {
		return testTemplate
	}
	return mainTemplate
}

func Generate(dot *Dot, name, out string) {
	helpers.Generate(templatesParsed, dot, name, out)
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

func consSlice(args []Arg) []string {
	a := make([]string, len(args))
	for i, arg := range args {
		switch arg.Type {
		case "string":
			a[i] = `""`
		case "int":
			a[i] = "int(0)"
		case "uint32":
			a[i] = "uint32(0)"
		case "*device.Peer":
			a[i] = "&device.Peer{}"
		case "error":
			a[i] = `errors.New("")`
		case "fmt.Stringer":
			a[i] = "testStringer{}"
		case "tai64n.Timestamp":
			a[i] = "tai64n.Timestamp{}"
		default:
			panic(fmt.Sprintf("invalid = name(%q):type(%q)", arg.Name, arg.Type))
		}
	}
	return a
}

func appendSlice(s1, s2 []string) []string {
	return append(append([]string{}, s1...), s2...)
}

func nilSlice(size int) []string {
	a := make([]string, size)
	for i := range a {
		a[i] = "nil"
	}
	return a
}

func joinStr(v []string) string {
	return strings.Join(v, ",")
}

func increment(n int) int {
	return add(n, 1)
}

func add(n1, n2 int) int {
	return n1 + n2
}

func mul(n1, n2 int) int {
	return n1 * n2
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
