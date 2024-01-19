package wgevents

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func Test_eventParser_ParseUDPGSODisabled(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		e      EventUDPGSODisabled
		wantOk bool
	}{
		{
			name: "invalid",
		},
		{
			name:  "valid",
			input: "disabled UDP GSO on 1.1.1.1, NIC(s) may not support checksum offload",
			e: EventUDPGSODisabled{
				OnLAddr: "1.1.1.1",
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &eventParser{}
			var e EventUDPGSODisabled
			if tt.wantOk {
				require.True(t, ev.ParseUDPGSODisabled(&e, tt.input))
				require.Equal(t, tt.e, e)
			} else {
				require.False(t, ev.ParseUDPGSODisabled(&e, tt.input))
			}
		})
	}

	t.Run("logger", func(t *testing.T) {
		ev := Events(func(ev Event) {
			// Test static fields
			require.Equal(t, NiceVerbosefUDPGSODisabled, ev.Nice())
			require.False(t, ev.IsErrorf())
			require.Equal(t, FormatVerbosefUDPGSODisabled, ev.Format())
			var buf bytes.Buffer
			// Test slogger
			ev.Slog(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
			require.NotEmpty(t, buf.Bytes())
			// Test args
			require.Equal(t, 1, len(ev.Args()))
			require.IsType(t, "", ev.Args()[0])
			require.Equal(t, "1.1.1.1", ev.Args()[0])
		})
		ev.Verbosef(fmt.Sprintf(FormatVerbosefUDPGSODisabled, "1.1.1.1"))
	})
}
