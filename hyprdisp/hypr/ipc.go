package hypr

import (
	"context"
	"fmt"
	"net"
	"path"
	"strings"
)

const (
	envHyprlandInstanceSignature = "HYPRLAND_INSTANCE_SIGNATURE"
	envXDGRuntimeDirectory       = "XDG_RUNTIME_DIR"
)

type HyprEvent struct {
	name string
	data string
}

func Listen(ctx context.Context) error {
	var (
		events chan HyprEvent = make(chan HyprEvent)
		errs   chan error     = make(chan error, 1)
	)
	defer func() {
		close(events)
		close(errs)
	}()

	go listenSocket(ctx, events, errs)

	for event := range events {
		if event.name != "monitoraddedv2" &&
			event.name != "monitorremoved" {
			fmt.Printf("Got irrelevant event: %v\n", event)
			continue
		}

		fmt.Printf("Got monitor event: %v\n", event)
	}

	for err := range errs {
		return err
	}

	return nil
}

func listenSocket(ctx context.Context, events chan HyprEvent, errs chan error) {
	var (
		socketPath string
		socketConn net.Conn
		err        error
	)

	socketPath, err = getEventsSocketPath()
	if err != nil {
		errs <- err
		return
	}

	socketConn, err = net.Dial("unix", socketPath)
	if err != nil {
		errs <- err
		return
	}

	var (
		bufferLength int    = 512
		bufferSize   int    = 0
		buffer       []byte = make([]byte, bufferLength)
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bufferSize, err = socketConn.Read(buffer)
		if err != nil {
			errs <- err
			return
		}

		if bufferSize <= 0 {
			continue
		}

		for _, event := range parseSocketEvents(string(buffer[:bufferSize])) {
			events <- event
		}
	}
}

func parseSocketEvents(data string) []HyprEvent {
	var (
		lines  []string    = strings.Split(data, "\n")
		events []HyprEvent = make([]HyprEvent, 0, len(lines))
	)

	for _, line := range lines {
		var parts []string

		if len(line) == 0 {
			continue
		}

		parts = strings.Split(line, ">>")
		events = append(events, HyprEvent{
			name: parts[0],
			data: parts[1],
		})

		// TODO identify how a line did not parse properly
		// TODO return the last line that did not parse. This is the "leftover"
	}

	return events
}

func getEventsSocketPath() (string, error) {
	var (
		runtimeDir string
		instanceID string
		err        error
	)

	runtimeDir, err = getRuntimePath()
	if err != nil {
		return "", nil
	}

	instanceID, err = getInstanceID()
	if err != nil {
		return "", nil
	}

	return path.Join(
		runtimeDir,
		instanceID,
		".socket2.sock",
	), nil
}

func getCommandsSocketPath() (string, error) {
	var (
		runtimeDir string
		instanceID string
		err        error
	)

	runtimeDir, err = getRuntimePath()
	if err != nil {
		return "", nil
	}

	instanceID, err = getInstanceID()
	if err != nil {
		return "", nil
	}

	return path.Join(
		runtimeDir,
		instanceID,
		".socket.sock",
	), nil
}
