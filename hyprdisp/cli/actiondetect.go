package cli

import (
	"context"
	"io/fs"
	"log/slog"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/profiles"
	"aeroheart.io/hyprdisp/sys"
)

type DetectAction struct {
	HyprLand hyprland.Service
	Profiles profiles.Service
}

func (action DetectAction) ID() string {
	return "detect"
}

func (action DetectAction) Execute(ctx context.Context) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var (
		monitors profiles.MonitorMap
		config   profiles.Config
	)
	monitors, err = action.Profiles.ConnectedMonitors(ctx)
	if err != nil {
		return err
	}

	logger.Info("Detecting current configuration")
	config, err = action.Profiles.Detect(ctx, monitors)
	switch err.(type) {
	case *fs.PathError:
		logger.Debug("Path error received from profile detection", slog.Any("error", err))
	default:
		return err

	}

	if !config.IsZero() {
		logger.Info("Configuration for monitors found. Exiting.")
		return nil
	}

	logger.Info("Configuration for monitors not found. Initializing.")
	_, err = action.Profiles.Init(ctx, monitors)
	return err
}
