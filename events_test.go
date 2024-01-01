package wgevents

import (
	"github.com/stretchr/testify/require"
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
}
