package hypr

import (
	"context"
	"net"
	"strings"
)

const (
	envHyprlandInstanceSignature = "HYPRLAND_INSTANCE_SIGNATURE"
	envXDGRuntimeDirectory       = "XDG_RUNTIME_DIR"
)

type Event struct {
	Name EventName
	Data []string
}

func StreamEvents(ctx context.Context, events chan Event, errs chan error) error {
	var (
		socketPath string
		socketConn net.Conn
		err        error
	)

	socketPath, err = getEventsSocketPath()
	if err != nil {
		return err
	}

	socketConn, err = net.Dial("unix", socketPath)
	if err != nil {
		return err
	}

	go watchEvents(
		ctx,
		socketConn,
		events,
		errs,
	)

	return nil
}

func watchEvents(ctx context.Context, socketConn net.Conn, events chan Event, errs chan error) {
	var (
		bufferLength int    = 512
		buffer       []byte = make([]byte, bufferLength)
		dataSize     int    = 0
		dataOverflow []byte
		err          error
	)

	defer func() {
		socketConn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		dataSize, err = socketConn.Read(buffer)
		if err != nil {
			errs <- err
			return
		}

		if dataSize <= 0 {
			continue
		}

		buffer = append(dataOverflow, buffer[:dataSize]...)

		var items []Event
		items, dataOverflow = parseEvents(buffer)

		for _, event := range items {
			events <- event
		}
	}
}

func parseEvents(data []byte) ([]Event, []byte) {
	var (
		lines     []string = strings.Split(string(data), "\n")
		lineCount int      = len(lines)
		lineFinal string   = lines[lineCount-1]
		events    []Event  = make([]Event, 0, lineCount)
	)

	if len(lineFinal) > 0 {
		lines = lines[:lineCount-1]
	}

	for _, line := range lines {
		var parts []string

		if len(line) == 0 {
			continue
		}

		parts = strings.Split(line, ">>")
		events = append(events, Event{
			Name: EventName(parts[0]),
			Data: strings.Split(parts[1], ","),
		})
	}

	return events, []byte(lineFinal)
}
