package hyprland

import (
	"context"
	"log"
	"os"
	"path"
	"strings"

	"aeroheart.io/hyprdisp/sys"
)

func (s defaultService) Apply(
	ctx context.Context,
	monitors []Monitor,
	workspaces []MonitorWorkspace,
) error {
	var (
		logger  *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		cfgPath string
		err     error
	)

	cfgPath, err = s.getConfigPath()
	if err != nil {
		return err
	}
	logger.Printf("Resovled configuration directory to: %s", cfgPath)

	logger.Printf("Writing monitors configuration to: %s", s.cfgMonitors)
	err = writeConfigMonitors(ctx, path.Join(cfgPath, s.cfgMonitors), monitors)
	if err != nil {
		return err
	}

	logger.Printf("Writing workspaces configuration to: %s", s.cfgWorkspaces)
	err = writeConfigWorkspaces(ctx, path.Join(cfgPath, s.cfgWorkspaces), workspaces)
	if err != nil {
		return err
	}

	return nil
}

func writeConfigMonitors(ctx context.Context, filepath string, monitors []Monitor) error {
	var (
		logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		lines  []string    = make([]string, 0, len(monitors))
	)

	logger.Printf("Marshalling monitor configurations")
	for _, monitor := range monitors {
		lines = append(lines, monitor.marshal())
	}
	lines = append(lines, "")

	logger.Printf("Saving monitor configurations")
	return writeConfig(filepath, []byte(strings.Join(lines, "\n")))
}

func writeConfigWorkspaces(ctx context.Context, filepath string, workspaces []MonitorWorkspace) error {
	var (
		logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		lines  []string    = make([]string, 0, len(workspaces))
	)

	logger.Printf("Marshalling workspace configurations")
	for _, workspace := range workspaces {
		lines = append(lines, workspace.marshal()...)
		lines = append(lines, "")
	}

	logger.Printf("Saving workspace configurations")
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
