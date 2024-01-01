package app

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"slices"
)

type Config struct {
	Imports []string `json:"import" yaml:"imports"`
	Events  map[string]struct {
		Type   YAMLEnum[LogType]  `json:"type" yaml:"type"`
		Level  YAMLEnum[LogLevel] `json:"level" yaml:"level"`
		Nice   string             `json:"nice" yaml:"nice"`
		Format string             `json:"format" yaml:"format"`
		Args   []*Arg             `json:"args" yaml:"args"`
		Custom bool               `json:"custom" yaml:"custom,omitempty"`
	} `json:"events" yaml:"events"`
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
	return errors.New("not a name: type pair")
}

func (a *Arg) MarshalYAML() (any, error) { return map[string]string{a.Name: a.Type}, nil }

type LogType string

const (
	LogTypeErrorf   LogType = "errorf"
	LogTypeVerbosef LogType = "verbosef"
)

func (LogType) Values() []string {
	return []string{
		string(LogTypeVerbosef),
		string(LogTypeErrorf),
	}
}

type LogLevel string

func (LogLevel) Values() []string {
	return []string{
		string(LogLevelError),
		string(LogLevelWarn),
		string(LogLevelInfo),
		string(LogLevelDebug),
	}
}

const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

type (
	YAMLEnum[E Enum] string
	Enum             interface{ Values() []string }
)

func (e *YAMLEnum[E]) UnmarshalYAML(n *yaml.Node) error {
	if err := n.Decode((*string)(e)); err != nil {
		return err
	} else if values := (*new(E)).Values(); !slices.Contains(values, string(*e)) {
		return fmt.Errorf("value of %q is invalid", *e)
	}
	return nil
}
