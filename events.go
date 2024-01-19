// Package wgevents is a logger able to turn format strings and args passed to the logger methods into concrete events.
// This is useful for structured logging as it allows arguments to be logged as separate fields.
// Support for [slog.Logger] is build in.
//
// This package is a separate module to pin this logger to a version of wireguard-go.
// This is needed since messages logged by wireguard-go may change without notice.
// Events can be added/removed/updated by editing `internal/events-generate/events.yml` and regenerating `events_generated.go`.
package wgevents

import "strings"

//go:generate go run github.com/point-c/wgevents/internal/cmd

// The init function also acts as an interface guard for [eventParser].
func init() { parser = new(eventParser) }

type eventParser struct{}

// ParseUDPGSODisabled is a method on eventParser that parses the EventUDPGSODisabled event.
// It takes an EventUDPGSODisabled instance, a format string, and variadic arguments.
// This method checks if the format string and input string conform to the expected structure of the event.
// If the format and input strings match, it extracts and sets the relevant information on the event instance and returns true.
// Otherwise, it returns false, indicating that the parsing was unsuccessful.
func (*eventParser) ParseUDPGSODisabled(ev *EventUDPGSODisabled, s string, _ ...any) (ok bool) {
	if prefix, suffix, ok := strings.Cut(ev.Format(), "%s"); ok && strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix) {
		ev.OnLAddr = strings.TrimPrefix(strings.TrimSuffix(s, suffix), prefix)
		return true
	}
	return false
}
