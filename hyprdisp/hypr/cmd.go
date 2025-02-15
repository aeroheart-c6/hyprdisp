package hypr

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Monitor struct {
	Num    int
	ID     string
	Make   string
	Model  string
	Serial string
}

func (m Monitor) Description() string {
	return fmt.Sprintf("%v %v %v",
		m.Make,
		m.Model,
		m.Serial,
	)
}

func GetMonitors() {
	var (
		socketPath string
		socketConn net.Conn
		err        error
	)

	socketPath, err = getCommandsSocketPath()
	if err != nil {
		return
	}

	socketConn, err = net.Dial("unix", socketPath)
	if err != nil {
		return
	}
	defer func() {
		socketConn.Close()
	}()

	socketConn.Write([]byte("monitors"))

	var (
		dataLength int    = 512
		dataCount  int    = 0
		data       []byte = make([]byte, dataLength)
		body       strings.Builder
	)

	for {
		socketConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		dataCount, err = socketConn.Read(data)
		if err != nil {
			fmt.Printf("%v", err)
			break
		}

		body.Write(data[:dataCount])

		if dataCount < dataLength {
			break
		}
	}

	fmt.Printf("MONITORS\n%v\n", string(body.String()))
}

func parseMonitorsBody(body string) []Monitor {
	return nil
}
