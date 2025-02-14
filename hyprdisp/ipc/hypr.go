package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
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

func ListenHyprEvents(ctx context.Context) (chan HyprEvent, error) {
	var (
		socketPath string
		err        error
	)

	socketPath, err = getHyprEventsSocketPath()
	if err != nil {
		return nil, err
	}

	var conn net.Conn

	conn, err = net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}

	var eventsChan chan HyprEvent = make(chan HyprEvent)
	go listen(ctx, conn, eventsChan)

	return eventsChan, nil
}

func listen(ctx context.Context, conn net.Conn, c chan HyprEvent) {
	var (
		data []byte = make([]byte, 512)
		size int
		err  error
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		size, err = conn.Read(data)
		if err != nil {
			fmt.Println("read error...")
		}

		if size <= 0 {
			continue
		}

		for _, event := range parse(data[:size]) {
			if event.name != "monitoraddedv2" &&
				event.name != "monitorremoved" {
				continue
			}
			c <- event
		}
	}
}

func parse(data []byte) []HyprEvent {
	var (
		lines  []string    = strings.Split(string(data), "\n")
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
	}

	return events
}

func getHyprEventsSocketPath() (string, error) {
	var (
		runtimeDir string
		instanceID string
		err        error
	)

	runtimeDir, err = getHyprRuntimePath()
	if err != nil {
		return "", nil
	}

	instanceID, err = getHyprInstanceID()
	if err != nil {
		return "", nil
	}

	return path.Join(
		runtimeDir,
		instanceID,
		".socket2.sock",
	), nil
}

func getHyprInstanceID() (string, error) {
	var (
		value string
		found bool
	)

	value, found = os.LookupEnv(envHyprlandInstanceSignature)
	if !found {
		return "", errors.New("hyprland instance signature not found")
	}

	return value, nil
}

func getHyprRuntimePath() (string, error) {
	var (
		value string
		found bool
	)

	value, found = os.LookupEnv(envXDGRuntimeDirectory)
	if !found {
		return "", errors.New("hyprland runtime path could not be acquired")
	}

	return path.Join(value, "hypr"), nil

}
