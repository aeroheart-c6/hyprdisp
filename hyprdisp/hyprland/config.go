package hyprland

import (
	"os"
	"path"
	"strings"
)

func Apply(monitors []Monitor, workspaces []MonitorWorkspace) error {
	var lines []string = make([]string, 0, len(monitors))

	// apply monitor configurations
	for _, monitor := range monitors {
		lines = append(lines, monitor.marshal())
	}

	lines = append(lines, "", "")

	// apply workspace configurations
	for _, workspace := range workspaces {
		lines = append(lines, workspace.marshal()...)
		lines = append(lines, "")
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

	filePath = path.Join(filePath, "actual.hypr-displays.conf")
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
