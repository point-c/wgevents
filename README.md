# wgevents

[![Go Reference](https://pkg.go.dev/badge/github.com/point-c/wgevents@v0.0.2.svg)](https://godocs.io/github.com/point-c/wgevents@v0.0.4)

## Overview

`wgevents` is a logging package for [`wireguard-go`](https://github.com/WireGuard/wireguard-go). It transforms format strings and args passed to logger methods into concrete events.

## Key Features

- **Event-based Logging:** Converts format strings and arguments into concrete logging events.
- **Structured Logging Support:** logs format arguments separate slog records. Integrates with `slog.Logger`.
- **Wireguard-go Version Pinning:** The package is designed as a separate module to pin to a specific version of `wireguard-go`. This ensures compatibility and stability even when `wireguard-go` messages change.

## Installation

To use wgevents in your Go project, install it using `go get`:

```bash
go get github.com/point-c/wgevents
```

## Usage

### Events Function

The `Events` function creates a new `device.Logger` capable of recognizing names and types of arguments passed to logging methods.

```go
func Events(fn func(Event)) *device.Logger
```

#### Parameters:
- `fn func(Event)`: A callback function that is called for each event logged.

#### Returns:
- `*device.Logger`: A logger compatible with `wireguard-go`.

### Event Handling

Each logged event is handled through the provided callback function. 
The package handles known events with specific types, and unknown events are categorized as `EventAny`.

## Example

### Basic Usage

```go
logger := wgevents.Events(func(e wgevents.Event) {
    // Handle logging event
})
// Use logger in your application
```

### Recognizing Events

A type switch can be used to recognize specific events:

```go
switch ev := ev.(type) {
case *wgevents.EventAssumingDefaultMTU:
case *wgevents.EventBindCloseFailed:
case *wgevents.EventBindUpdated:
// And so on...
}
```

### Logging

```go
logger := wgevents.Events(func(e wgevents.Event) {
    e.Slog(slog.Default())
})
```

## Contributing

This package allows for updates and additions of events. Modify `internal/events-generate/events.yml` and regenerate `events_generated.go` to update event handling.

### File Layout

Each section under the `events` key corresponds to a specific event that can be logged:

```yaml
<event name>:
    type: string
    level: string
    nice: string (optional)
    format: string
    args: list (optional)
    custom: bool (optional)
```

### Key Definitions

- `type`: The type of logging or messaging method to be used. Only `errorf` for error logging, and `verbosef` for verbose logging are valid.
- `level`: The level at which the logging should occur. Valid values are `debug`, `info`, `warn`, and `error`.
- `nice`: An optional human-readable message for the event. If not specified `format` will be used.
- `format`: The format string used by the `wireguard-go` library.
- `args`: The arguments used in the format string in the following format:
  
    ```yaml
    - <name>: <type>
    - <name>: <type>
    ```

- `custom`: An optional flag indicating that parsing this event requires custom logic to be implemented.

### Imports

The final section of the configuration file lists the external packages used as part of the event processing mechanism. These are mentioned under the `imports` key.

## Testing

The package includes tests that demonstrate its functionality. Use Go's testing tools to run the tests:

```bash
go test ./...
```

## Godocs

To regenerate godocs:

```bash
go generate -tags docs ./...
```

## Code Generation

To regenerate generated packages:

```bash
go generate ./...
```