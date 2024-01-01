package helpers

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func preserveExit(t *testing.T) func() {
	t.Helper()
	original := exit
	return func() { t.Helper(); exit = original }
}

func TestCheck(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		defer preserveExit(t)
		var run bool
		exit = func(code int) { run = true; require.Equal(t, code, 1) }
		Check(errors.New("test"))
		require.True(t, run)
	})

	t.Run("no error", func(t *testing.T) {
		defer preserveExit(t)
		exit = func(int) { t.Fail() }
		Check(nil)
	})
}

func TestMust(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		defer preserveExit(t)
		var run bool
		exit = func(code int) { run = true; require.Equal(t, code, 1) }
		require.Equal(t, 1, Must(1, errors.New("test")))
		require.True(t, run)
	})

	t.Run("no error", func(t *testing.T) {
		defer preserveExit(t)
		exit = func(int) { t.Fail() }
		require.Equal(t, 1, Must(1, nil))
	})
}
