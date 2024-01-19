// Package templates manages the generation of code based on embedded templates.
// It leverages Go's text/template package for template processing and embeds the templates using Go's embed package.
// The package is designed to facilitate the generation of code for events for structured logging.
package templates

import (
	"embed"
	"fmt"
	helpers "github.com/point-c/generator-helpers"
	"strings"
	"text/template"
)

const (
	// mainTemplate is the entrypoint for the main code.
	mainTemplate = "events.gotmpl"
	// testTemplate is the entrypoint for the test code.
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

// GetTemplate gets the name of the template to use. An argument of false will return the name of the main template.
func GetTemplate(test bool) string {
	if test {
		return testTemplate
	}
	return mainTemplate
}

// Generate processes the provided template with the given data and outputs to the specified location.
func Generate(dot *Dot, name, out string) {
	helpers.Generate(templatesParsed, dot, name, out)
}

// isError formats the template context to include an error flag and a list of events.
func isError(ev []*Event, isErr bool) any {
	return &struct {
		IsError bool
		Events  []*Event
	}{
		IsError: isErr,
		Events:  ev,
	}
}

// consSlice constructs a slice of initializers based on the types of the provided arguments.
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

// appendSlice appends two slices of strings.
func appendSlice(s1, s2 []string) []string {
	return append(append([]string{}, s1...), s2...)
}

// nilSlice creates a slice of strings, all set to "nil".
func nilSlice(size int) []string {
	a := make([]string, size)
	for i := range a {
		a[i] = "nil"
	}
	return a
}

// joinStr joins a slice of strings using a comma separator.
func joinStr(v []string) string {
	return strings.Join(v, ",")
}

// increment increases a number by one.
func increment(n int) int {
	return add(n, 1)
}

// add adds two numbers.
func add(n1, n2 int) int {
	return n1 + n2
}

// mul multiplies two numbers.
func mul(n1, n2 int) int {
	return n1 * n2
}

type (
	Dot struct {
		Package  string   // Package name
		Imports  []string // Extra imports
		Events   []*Event // All events
		Verbosef []*Event // Verbose events
		Errorf   []*Event // Error events
		Custom   []*Event // Events with custom parsers
	}
	Event struct {
		Name   string // Camel case name with no spaces
		Type   string // "errorf" or "verbosef"
		Level  string // Slog levels of "debug", "info", "error", or "warn"
		Nice   string // A non format string version of the message to print
		Format string // The format string from the wireguard-go library to match
		Args   []Arg  // Ordered list of args used in the format string
		Custom bool   // Whether the event requires a custom parser
	}
	Arg struct {
		Name string // Name to use for the argument
		Type string // Type of the argument
	}
)
