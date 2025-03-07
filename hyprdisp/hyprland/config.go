package hyprland

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func Apply(monitors []Monitor) error {
	var lines []string = make([]string, 0, len(monitors))

	for _, monitor := range monitors {
		lines = append(lines, serializeToConfig(monitor))
	}

	var (
		filePath string
		file     *os.File
		err      error
	)
	filePath, err = getConfigPath()
	if err != nil {
		return err
	}

	filePath = path.Join(filePath, "hypr-displays.conf")
	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return err
	}

	return nil
}

func serializeToConfig(monitor Monitor) string {
	return fmt.Sprintf("monitor = %s, %s, %s, %s",
		monitor.Name,
		monitor.Resolution,
		monitor.Position,
		monitor.Scale,
	)
}
