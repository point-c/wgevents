package helpers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGoFmt(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		input := "package main\nimport \"fmt\"\nfunc main() {\nfmt.Println(`Hello World`)}"
		expected := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(`Hello World`)\n}\n"
		result, err := GoFmt([]byte(input))
		require.NoError(t, err)
		require.Equal(t, expected, string(result))
	})
	t.Run("error", func(t *testing.T) {
		input := "package mainimport fmtfunc main() {fmt.Println(`Hello World`)}"
		result, err := GoFmt([]byte(input))
		require.Error(t, err)
		require.Equal(t, string(result), input)
	})
}
