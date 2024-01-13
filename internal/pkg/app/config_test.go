package app

import (
	helpers "github.com/point-c/generator-helpers"
	"github.com/point-c/wgevents/internal/pkg/templates"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArg_MarshalYAML(t *testing.T) {
	arg := Arg{
		Name: "foo",
		Type: "bar",
	}
	v, err := arg.MarshalYAML()
	require.NoError(t, err)
	require.Equal(t, map[string]string{"foo": "bar"}, v)
}

func TestArg_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		expected Arg
		input    string
		wantErr  bool
	}{
		{
			name:    "invalid",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			input:   `"foo"`,
			wantErr: true,
		},
		{
			name: "valid",
			expected: Arg{
				Name: "foo",
				Type: "bar",
			},
			input: "foo: bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(tt.input), &n))
			var a Arg
			if err := a.UnmarshalYAML(&n); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, a)
			}
		})
	}
}

func TestArg_unmarshalArg(t *testing.T) {
	tests := []struct {
		name    string
		want    Arg
		input   map[string]string
		wantErr bool
	}{
		{
			name: "valid",
			want: Arg{
				Name: "foo",
				Type: "bar",
			},
			input: map[string]string{"foo": "bar"},
		},
		{
			name:    "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Arg
			if err := a.unmarshalArg(tt.input); tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, a)
			}
		})
	}
}

func TestLogLevel_Values(t *testing.T) {
	var l LogLevel
	require.NotEmpty(t, l.Values())
}

func TestLogType_Values(t *testing.T) {
	var l LogType
	require.NotEmpty(t, l.Values())
}

type testYamlEnum struct{}

func (testYamlEnum) Values() []string { return []string{"test"} }

func TestYAMLEnum_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		e       string
		input   string
		wantErr bool
	}{
		{
			name:    "invalid yaml",
			input:   "- foo",
			wantErr: true,
		},
		{
			name:    "invalid enum",
			input:   "foo",
			wantErr: true,
		},
		{
			name:  "valid",
			e:     "test",
			input: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var n yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(tt.input), &n))
			var a YAMLEnum[testYamlEnum]
			if tt.wantErr {
				require.Error(t, a.UnmarshalYAML(&n))
			} else {
				require.NoError(t, a.UnmarshalYAML(&n))
				require.Equal(t, YAMLEnum[testYamlEnum](tt.e), a)
			}
		})
	}
}

func TestCfg2Dot(t *testing.T) {
	firstEv := templates.Event{
		Name:   "first",
		Type:   "errorf",
		Level:  "error",
		Nice:   "foo",
		Format: "%s",
		Args:   []templates.Arg{{Name: "foo", Type: "bar"}},
		Custom: true,
	}
	secondEv := templates.Event{
		Name:   "second",
		Type:   "verbosef",
		Level:  "info",
		Nice:   "foo",
		Format: "%s",
		Args:   []templates.Arg{{Name: "foo", Type: "bar"}},
	}
	cfgFn := filepath.Join(t.TempDir(), "cfg.yml")
	require.NoError(t, os.WriteFile(cfgFn, []byte(strings.ReplaceAll(`imports:
	- bar
	- baz
events:
	first:
		custom: true
		nice: foo
		format: "%s"
		type: errorf
		level: error
		args:
			- foo: bar
	second:
		nice: foo
		level: info
		type: verbosef
		format: "%s"
		args:
			- foo: bar
`, "\t", "    ")), os.ModePerm))
	cfg, err := helpers.UnmarshalYAML[Config](cfgFn)
	require.NoError(t, err)
	d := Cfg2Dot(&cfg, "foo")
	require.Equal(t, "foo", d.Package)
	require.Equal(t, []string{"bar", "baz"}, d.Imports)
	require.Len(t, d.Events, 2)
	require.Equal(t, firstEv, *d.Events[0])
	require.Equal(t, secondEv, *d.Events[1])
	require.Len(t, d.Verbosef, 1)
	require.Equal(t, secondEv, *d.Verbosef[0])
	require.Len(t, d.Errorf, 1)
	require.Equal(t, firstEv, *d.Errorf[0])
	require.Len(t, d.Custom, 1)
	require.Equal(t, firstEv, *d.Custom[0])
}
