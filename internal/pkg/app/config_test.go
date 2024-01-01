package app

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
