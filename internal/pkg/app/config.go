package app

import (
	"errors"
	"fmt"
	helpers "github.com/point-c/generator-helpers"
	"github.com/point-c/wgevents/internal/pkg/templates"
	_ "golang.zx2c4.com/wireguard/device"
	"gopkg.in/yaml.v3"
	"slices"
)

// init function sets default values for the environment variables used in the helper package.
func init() {
	helpers.EnvDefaultGoPackage = "wgevents"
	helpers.EnvDefaultGoFile = "events.go"
}

// Config struct defines the main configuration for generating the code.
type Config struct {
	// Additional imports for the generated code.
	// For the main code `slog.Logger` is already included.
	// For test code `github.com/stretchr/testify/require`, `testing`, `bytes`, `errors`, and `slog.Logger` are already included.
	Imports []string `json:"import" yaml:"imports"`
	Events  map[string]struct {
		Type   YAMLEnum[LogType]  `json:"type" yaml:"type"`               // Event type (errorf, verbosef).
		Level  YAMLEnum[LogLevel] `json:"level" yaml:"level"`             // Log level (error, warn, info, debug).
		Nice   string             `json:"nice" yaml:"nice"`               // User-friendly message.
		Format string             `json:"format" yaml:"format"`           // Format string for identifying events.
		Args   []*Arg             `json:"args" yaml:"args"`               // Arguments for the format string.
		Custom bool               `json:"custom" yaml:"custom,omitempty"` // Flag for custom parsing.
	} `json:"events" yaml:"events"`
}

var _ interface {
	yaml.Marshaler
	yaml.Unmarshaler
} = (*Arg)(nil)

// Arg represents an argument in the log format string.
// An Arg is represented in YAML as a map with a single element.
type Arg struct {
	Name string
	Type string
}

// UnmarshalYAML custom unmarshals YAML node into Arg struct.
func (a *Arg) UnmarshalYAML(n *yaml.Node) error {
	var m map[string]string
	if err := n.Decode(&m); err != nil {
		return err
	}
	return a.unmarshalArg(m)
}

// unmarshalArg processes a map and assigns it to the Arg struct.
func (a *Arg) unmarshalArg(m map[string]string) error {
	// More than one element is invalid
	if len(m) == 1 {
		// Get the first element and return.
		// Since the key is unknown use a loop.
		for n, t := range m {
			a.Name, a.Type = n, t
			return nil
		}
	}
	return errors.New("not a name: type pair")
}

func (a *Arg) MarshalYAML() (any, error) { return map[string]string{a.Name: a.Type}, nil }

// LogType are the valid log levels available from [device.Logger].
type LogType string

const (
	LogTypeErrorf   LogType = "errorf"
	LogTypeVerbosef LogType = "verbosef"
)

// Values method returns all valid values for LogType.
func (LogType) Values() []string {
	return []string{
		string(LogTypeVerbosef),
		string(LogTypeErrorf),
	}
}

// LogLevel are the valid log levels available from [slog.Logger].
type LogLevel string

// Values method returns all valid values for LogLevel.
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
	// YAMLEnum handles unmarshalling and validation of enums from yaml.
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

// Cfg2Dot converts the configuration into a templates.Dot for code generation.
func Cfg2Dot(cfg *Config, pkg string) *templates.Dot {
	d := templates.Dot{
		Package: pkg,
		Imports: cfg.Imports,
	}

	// Process and filter events.
	for name, ev := range cfg.Events {
		// Convert struct
		ev := templates.Event{
			Name:   name,
			Type:   string(ev.Type),
			Level:  string(ev.Level),
			Nice:   helpers.IfStringEmptyUseDefault(ev.Nice, ev.Format),
			Format: ev.Format,
			// Convert [Arg]s to [tmeplates.Arg]s
			Args: func(args []templates.Arg) []templates.Arg {
				for i, arg := range ev.Args {
					args[i].Name = arg.Name
					args[i].Type = arg.Type
				}
				return args
			}(make([]templates.Arg, len(ev.Args))), // Preallocate slice
			Custom: ev.Custom,
		}

		// Add to events
		d.Events = append(d.Events, &ev)

		// Add to custom
		if ev.Custom {
			d.Custom = append(d.Custom, &ev)
		}

		// Filter into correct log type slice
		switch LogType(ev.Type) {
		case LogTypeErrorf:
			d.Errorf = append(d.Errorf, &ev)
		case LogTypeVerbosef:
			d.Verbosef = append(d.Verbosef, &ev)
		}
	}
	return &d
}
