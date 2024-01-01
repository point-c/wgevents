package app

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOutputFilename(t *testing.T) {
	for _, c := range []struct {
		name, input, expected string
	}{
		{
			name:     "with ext",
			input:    "test.go",
			expected: "test_generated.go",
		},
		{
			name:     "without ext",
			input:    "test",
			expected: "test_generated.go",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.expected, OutputFilename(c.input))
		})
	}
}

func TestEnvOrDefault(t *testing.T) {
	type args struct {
		key string
		def string
	}
	type test struct {
		name string
		args args
		want string
		env  map[string]string
	}

	tests := []test{
		{
			name: "default",
			args: args{
				key: uuid.New().String(),
				def: "test",
			},
			want: "test",
		},
		func(key string) test {
			return test{
				name: "env",
				args: args{key: key},
				want: "test",
				env:  map[string]string{key: "test"},
			}
		}(uuid.New().String()),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != nil {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
			}

			require.Equal(t, tt.want, EnvOrDefault(tt.args.key, tt.args.def))
		})
	}
}

func TestIgnoreForCoverage(t *testing.T) {
	// These will only return default unless called by go generate
	_ = GoFile()
	_ = GoPackage()
}

func TestIfStringEmptyUseDefault(t *testing.T) {
	type args struct {
		s   string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "non-empty",
			args: args{
				s:   "foo",
				def: "bar",
			},
			want: "foo",
		},
		{
			name: "empty",
			args: args{
				s:   "",
				def: "bar",
			},
			want: "bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IfStringEmptyUseDefault(tt.args.s, tt.args.def); got != tt.want {
				t.Errorf("IfStringEmptyUseDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
