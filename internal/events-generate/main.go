package main

import (
	"bytes"
	"embed"
	"errors"
	"go/format"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

//go:embed events.yml
var eventsYML []byte

//go:embed templates/*
var templates embed.FS

const (
	templatesSubdir = "templates"
	mainTemplate    = "events.gotmpl"
)

var outputFilename = strings.TrimSuffix(os.Getenv("GOFILE"), ".go") + "_generate.go"

var extraImports = []string{
	"log/slog",
}

func init() {
	if os.Getenv("GOPACKAGE") == "" {
		os.Setenv("GOPACKAGE", "wgevents")
		outputFilename = filepath.Join("pkg", "wg", "wglog", "wgevents", "events"+outputFilename)
	}
}

func main() {
	var buf bytes.Buffer
	check(getTemplates().ExecuteTemplate(&buf, mainTemplate, NewDot(unmarshalEvents())))
	check(os.WriteFile(outputFilename, ignore(gofmt(buf.Bytes())), os.ModePerm))
}

func unmarshalEvents() *Config {
	var events Config
	check(yaml.Unmarshal(eventsYML, &events))
	return &events
}

func getTemplates() *template.Template {
	return template.Must(template.New("").Funcs(template.FuncMap{
		"isError": IsError,
	}).ParseFS(must(fs.Sub(templates, templatesSubdir)), "*"))
}

func gofmt(src []byte) (dst []byte, err error) {
	dst, err = format.Source(src)
	if err != nil {
		slog.Error("failed to format code", "error", err)
		dst = src
	}
	return
}

func ignore[T any](t T, _ error) T {
	return t
}

func must[T any](t T, err error) T {
	if err != nil {
		slog.Error("failed to get value", "error", err)
		os.Exit(1)
	}
	return t
}

func check(err error) {
	if err != nil {
		slog.Error("operation failed", "error", err)
		os.Exit(1)
	}
}

type Config struct {
	Imports Imports `json:"import" yaml:"imports"`
	Events  Events  `json:"events" yaml:"events"`
}

type Imports []string

func (i *Imports) UnmarshalYAML(n *yaml.Node) error {
	if err := n.Decode((*[]string)(i)); err != nil {
		return err
	}
	*i = append(*i, extraImports...)
	slices.Sort(*i)
	return nil
}

type Events []*Event

func (e *Events) UnmarshalYAML(n *yaml.Node) error {
	var m map[string]*Event
	if err := n.Decode(&m); err != nil {
		return err
	}
	*e = make([]*Event, 0, len(m))
	for name, value := range m {
		value.Name = name
		if value.Nice == "" {
			value.Nice = value.Format
		}
		*e = append(*e, value)
	}
	slices.SortFunc(*e, func(a, b *Event) int { return strings.Compare(a.Name, b.Name) })
	return nil
}

type Dot struct {
	Imports  []string
	Events   []*Event
	Verbosef []*Event
	Errorf   []*Event
	Custom   []*Event
}

func (*Dot) Package() string { return os.Getenv("GOPACKAGE") }

func NewDot(cfg *Config) *Dot {
	d := Dot{
		Imports: cfg.Imports,
		Events:  cfg.Events,
	}

	for _, ev := range d.Events {
		d.filterType(ev)
		d.filterCustom(ev)
	}
	return &d
}

func (d *Dot) filterCustom(ev *Event) {
	if ev.Custom {
		d.Custom = append(d.Custom, ev)
	}
}

func (d *Dot) filterType(ev *Event) {
	switch ev.Type {
	case LogTypeErrorf:
		d.Errorf = append(d.Errorf, ev)
	case LogTypeVerbosef:
		d.Verbosef = append(d.Verbosef, ev)
	default:
		panic("invalid type: " + ev.Type)
	}
}

type Event struct {
	Name   string   `json:"-" yaml:"-"`
	Type   LogType  `json:"type" yaml:"type"`
	Level  LogLevel `json:"level" yaml:"level"`
	Nice   string   `json:"nice" yaml:"nice"`
	Format string   `json:"format" yaml:"format"`
	Args   []*Arg   `json:"args" yaml:"args"`
	Custom bool     `json:"custom" yaml:"custom,omitempty"`
}

var _ interface {
	yaml.Marshaler
	yaml.Unmarshaler
} = (*Arg)(nil)

type Arg struct {
	Name string
	Type string
}

func (a *Arg) UnmarshalYAML(n *yaml.Node) error {
	var m map[string]string
	if err := n.Decode(&m); err != nil {
		return err
	}
	return a.unmarshalArg(m)
}

func (a *Arg) unmarshalArg(m map[string]string) error {
	for n, t := range m {
		a.Name, a.Type = n, t
		return nil
	}
	return errors.New("no a name: type pair")
}

func (a *Arg) MarshalYAML() (any, error) { return map[string]string{a.Name: a.Type}, nil }

type LogType string

const (
	LogTypeErrorf   LogType = "errorf"
	LogTypeVerbosef LogType = "verbosef"
)

type LogLevel string

const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

func IsError(ev []*Event, isErr bool) any {
	return &struct {
		IsError bool
		Events  []*Event
	}{
		IsError: isErr,
		Events:  ev,
	}
}
