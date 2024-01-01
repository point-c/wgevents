package helpers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalYAML(t *testing.T) {
	type TestYAMLStruct struct {
		Foo string `yaml:"foo"`
		Baz int    `yaml:"baz"`
	}

	for _, c := range []struct {
		name     string
		filename string
		err      bool
	}{
		{
			name:     "valid",
			filename: "test.yml",
		},
		{
			name:     "invalid",
			filename: "invalid.yml",
			err:      true,
		},
		{
			name:     "file not exist",
			filename: "NOT_EXISTS.yml",
			err:      true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			v, err := UnmarshalYAML[TestYAMLStruct](c.filename)
			if c.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, TestYAMLStruct{
					Foo: "bar",
					Baz: 123,
				}, v)
			}
		})
	}
}
