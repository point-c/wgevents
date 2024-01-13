package templates

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"log/slog"
	"path/filepath"
	"testing"
)

func TestIsError(t *testing.T) {
	require.NotNil(t, isError(nil, true))
}

func TestGetTemplate(t *testing.T) {
	require.Equal(t, mainTemplate, GetTemplate(false))
	require.Equal(t, testTemplate, GetTemplate(true))
}

func TestAdd(t *testing.T) {
	require.Equal(t, 3, add(1, 2))
}

func TestMul(t *testing.T) {
	require.Equal(t, 6, mul(2, 3))
}

func TestIncrement(t *testing.T) {
	require.Equal(t, 4, increment(3))
}

func TestJoinStr(t *testing.T) {
	require.Equal(t, "foo,bar", joinStr([]string{"foo", "bar"}))
}

func TestNilSlice(t *testing.T) {
	require.Equal(t, []string{"nil", "nil"}, nilSlice(2))
}

func TestAppendSlice(t *testing.T) {
	require.Equal(t, []string{"foo", "bar", "baz", "zab"}, appendSlice([]string{"foo", "bar"}, []string{"baz", "zab"}))
}

func TestConsSlice(t *testing.T) {
	require.Panics(t, func() {
		consSlice([]Arg{{}})
	})
	require.Equal(t, []string{
		`""`,
		"int(0)",
		"uint32(0)",
		"&device.Peer{}",
		`errors.New("")`,
		"testStringer{}",
		"tai64n.Timestamp{}",
	}, consSlice([]Arg{
		{Type: "string"},
		{Type: "int"},
		{Type: "uint32"},
		{Type: "*device.Peer"},
		{Type: "error"},
		{Type: "fmt.Stringer"},
		{Type: "tai64n.Timestamp"},
	}))
}

func TestGenerate(t *testing.T) {
	defer slog.SetDefault(slog.Default())
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	out := filepath.Join(t.TempDir(), "out.go")
	Generate(&Dot{Package: "foo"}, mainTemplate, out)
	require.Empty(t, buf.Bytes())
	require.FileExists(t, out)
}
