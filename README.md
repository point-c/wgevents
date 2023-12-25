# wgevents

## Overview

`wgevents` is a logging package for [`wireguard-go`](https://github.com/WireGuard/wireguard-go). It transforms format strings and args passed to logger methods into concrete events.

## Key Features

- **Event-based Logging:** Converts format strings and arguments into concrete logging events.
- **Structured Logging Support:** logs format arguments separate slog records. Integrates with `slog.Logger`.
- **Wireguard-go Version Pinning:** The package is designed as a separate module to pin to a specific version of `wireguard-go`. This ensures compatibility and stability even when `wireguard-go` messages change.

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

```go
package main

import (
    "github.com/point-c/wgevents"
)

func main() {
    logger := wgevents.Events(func(e wgevents.Event) {
        // Handle logging event
    })
    // Use logger in your application
}
```

A type switch can be used to recognize specific events:

```go
switch ev := ev.(type) {
case *wgevents.EventAssumingDefaultMTU:
case *wgevents.EventBindCloseFailed:
case *wgevents.EventBindUpdated:
// And so on...
}
```

## Contributing

This package allows for updates and additions of events. Modify `internal/events-generate/events.yml` and regenerate `events_generated.go` to update event handling.
