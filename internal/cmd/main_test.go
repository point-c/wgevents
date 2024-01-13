package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_main(t *testing.T) {
	defer func(args []string) { os.Args = args }(os.Args)
	tmpDir := t.TempDir()
	events := filepath.Join(tmpDir, "events.yml")
	require.NoError(t, os.WriteFile(events, []byte(`events:
`), os.ModePerm))
	out := filepath.Join(tmpDir, "out_generated.go")
	tout := filepath.Join(tmpDir, "tout_generated.go")
	os.Args = []string{"test", "-out", out, "-tout", tout, "-config", events}
	main()
	require.FileExists(t, out)
	require.FileExists(t, tout)
}
