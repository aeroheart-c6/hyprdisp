package hyprland

import (
	"os"
	"path"
	"strings"
)

func (s defaultService) Apply(monitors []Monitor, workspaces []MonitorWorkspace) error {
	var (
		basepath string
		err      error
	)

	basepath, err = s.getConfigPath()
	if err != nil {
		return err
	}

	err = writeConfigMonitors(
		path.Join(basepath, "hypr-monitors.conf"),
		monitors,
	)
	if err != nil {
		return err
	}

	err = writeConfigWorkspaces(
		path.Join(basepath, "hypr-workspaces.conf"),
		workspaces,
	)
	if err != nil {
		return err
	}

	return nil
}

func writeConfigMonitors(filepath string, monitors []Monitor) error {
	var lines []string = make([]string, 0, len(monitors))

	// apply monitor configurations
	for _, monitor := range monitors {
		lines = append(lines, monitor.marshal())
	}

	lines = append(lines, "")

	return writeConfig(filepath, []byte(strings.Join(lines, "\n")))
}

func writeConfigWorkspaces(filepath string, workspaces []MonitorWorkspace) error {
	var lines []string = make([]string, 0, len(workspaces))

	// apply workspace configurations
	for _, workspace := range workspaces {
		lines = append(lines, workspace.marshal()...)
		lines = append(lines, "")
	}

	return writeConfig(filepath, []byte(strings.Join(lines, "\n")))
}

func writeConfig(filepath string, data []byte) error {
	var (
		file *os.File
		err  error
	)
	file, err = os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
