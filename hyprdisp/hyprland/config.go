package hyprland

import (
	"context"
	"log/slog"
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
		logger  *slog.Logger
		cfgPath string
		err     error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	cfgPath, err = s.getConfigPath()
	if err != nil {
		return err
	}
	logger.Info("Resolved configuration directory", slog.String("path", cfgPath))

	logger.Info("Writing monitors configuration", slog.String("file", s.cfgMonitors))
	err = writeConfigMonitors(ctx, path.Join(cfgPath, s.cfgMonitors), monitors)
	if err != nil {
		return err
	}

	logger.Info("Writing workspaces configuration", slog.String("file", s.cfgWorkspaces))
	err = writeConfigWorkspaces(ctx, path.Join(cfgPath, s.cfgWorkspaces), workspaces)
	if err != nil {
		return err
	}

	return nil
}

func writeConfigMonitors(ctx context.Context, filepath string, monitors []Monitor) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	logger.Info("Marshalling monitor configurations")

	var lines []string = make([]string, 0, len(monitors))
	for _, monitor := range monitors {
		lines = append(lines, monitor.marshal())
	}
	lines = append(lines, "")

	logger.Info("Saving monitor configurations")
	return writeConfig(filepath, []byte(strings.Join(lines, "\n")))
}

func writeConfigWorkspaces(ctx context.Context, filepath string, workspaces []MonitorWorkspace) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	logger.Info("Marshalling workspace configurations")
	var lines []string = make([]string, 0, len(workspaces))

	for _, workspace := range workspaces {
		lines = append(lines, workspace.marshal()...)
		lines = append(lines, "")
	}

	logger.Info("Saving workspace configurations")
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
