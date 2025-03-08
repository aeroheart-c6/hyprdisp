package hyprland

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func (s defaultService) GetMonitors() ([]Monitor, error) {
	var (
		socketPath string
		socketConn net.Conn
		err        error
	)

	socketPath, err = getCommandsSocketPath()
	if err != nil {
		return nil, err
	}

	socketConn, err = net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		socketConn.Close()
	}()

	socketConn.Write([]byte("monitors"))

	var (
		bufferSize int    = 512
		buffer     []byte = make([]byte, bufferSize)
		dataSize   int    = 0
		data       strings.Builder
	)

	for {
		socketConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		dataSize, err = socketConn.Read(buffer)
		if err != nil {
			fmt.Printf("%v", err)
			break
		}

		data.Write(buffer[:dataSize])

		if dataSize < bufferSize {
			break
		}
	}

	return parseMonitorsPayload(data.String())
}

var (
	regexMonitorID  *regexp.Regexp = regexp.MustCompile(`Monitor (?P<id>[\w-]+) \(ID (?P<no>\d+)\)`)
	regexMonitorRes *regexp.Regexp = regexp.MustCompile(`` +
		`(?P<resolution>\d+x\d+@\d+\.\d+)` +
		` at ` +
		`(?P<position>\d+x\d+)`,
	)
)

func parseMonitorsPayload(body string) ([]Monitor, error) {
	var (
		lines        []string  = strings.Split(body, "\n")
		monitors     []Monitor = make([]Monitor, 0, 5)
		bodyStartIdx int       = 0
		bodyEndIdx   int       = len(lines)
	)

	for {
		if bodyStartIdx > bodyEndIdx {
			break
		}

		var (
			monitor Monitor
			lineIdx int
			line    string
		)
		for lineIdx, line = range lines[bodyStartIdx:bodyEndIdx] {
			line = strings.TrimSpace(line)

			if len(line) == 0 {
				// Terminating line
				if !monitor.IsZero() {
					monitors = append(monitors, monitor)
				}
				break
			} else if lineIdx == 0 {
				// Monitor IDs
				var matches []string = regexMonitorID.FindStringSubmatch(line)
				if matches == nil {
					return nil, errors.New("invalid monitor payload")
				}

				monitor = Monitor{}
				monitor.set("id", matches[2])
				monitor.set("name", matches[1])
			} else if lineIdx == 1 {
				// Monitor Resolution
				var matches []string = regexMonitorRes.FindStringSubmatch(line)
				if matches == nil {
					return nil, errors.New("invalid monitor resolution spec")
				}

				monitor.set("resolution", matches[1])
				monitor.set("position", matches[2])
			} else {
				var (
					field string
					value string
					found bool
					err   error
				)

				field, value, found = strings.Cut(line, ":")
				if !found {
					return nil, errors.New("uanble to parse a field / value pair")
				}

				err = monitor.set(field, strings.TrimSpace(value))
				if err != nil {
					fmt.Printf("whoops: %v\n", err)
				}
			}
		}

		bodyStartIdx = bodyStartIdx + lineIdx + 1
	}

	return monitors, nil
}
