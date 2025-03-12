package hyprland

import (
	"context"
	"net"
	"strings"
)

func (s defaultService) StreamEvents(ctx context.Context) (chan Event, chan error, error) {
	var (
		socketPath string
		socketConn net.Conn
		err        error
	)

	socketPath, err = getEventsSocketPath()
	if err != nil {
		return nil, nil, err
	}

	socketConn, err = net.Dial("unix", socketPath)
	if err != nil {
		return nil, nil, err
	}

	var (
		events chan Event = make(chan Event)
		errs   chan error = make(chan error, 1)
	)

	go watchEvents(
		ctx,
		socketConn,
		events,
		errs,
	)

	return events, errs, nil
}

func watchEvents(ctx context.Context, socketConn net.Conn, events chan Event, errs chan error) {
	var (
		dataLength int    = 512
		data       []byte = make([]byte, dataLength)
		dataAvail  int    = 0
		dataExcess []byte
		err        error
	)

	defer func() {
		socketConn.Close()
		close(events)
		close(errs)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		dataAvail, err = socketConn.Read(data)
		if err != nil {
			errs <- err
			return
		}

		if dataAvail <= 0 {
			continue
		}

		data = append(dataExcess, data[:dataAvail]...)

		var items []Event
		items, dataExcess = parseEvents(data)

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
