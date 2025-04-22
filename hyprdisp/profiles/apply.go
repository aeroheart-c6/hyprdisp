package profiles

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/hyprpanel"
	"aeroheart.io/hyprdisp/sys"
)

const (
	keyDefaultPanelMain string = "main"
	keyDefaultPanelSub  string = "sub"
)

func (s defaultService) Apply(ctx context.Context, config Config) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var (
		origPath string
		confPath string
	)

	logger.Info("Applying profile: Panels", slog.String("id", config.ID))
	err = s.applyPanels(ctx, config)
	if err != nil {
		logger.Info("Unable to apply panel configuration", slog.Any("error", err))
	}

	logger.Info("Applying profile: Monitors and Workspaces", slog.String("id", config.ID))
	err = s.applyMonitors(ctx, config.Monitors)
	if err != nil {
		return err // TODO should probably try to roll back???
	}

	logger.Info("Applying profile: Setting current profile", slog.String("id", config.ID))
	confPath, err = s.getConfigPath()
	if err != nil {
		return err
	}
	confPath = path.Join(confPath, fmt.Sprintf("%s.current.toml", config.ID))
	origPath = fmt.Sprintf("%s.toml", config.ID)

	_, err = os.Lstat(confPath)
	if err == nil {
		err = os.Remove(confPath)
		if err != nil {
			return err
		}
	} else {
		var (
			pathErr *fs.PathError
			ok      bool
		)

		pathErr, ok = err.(*fs.PathError)
		if !ok || !errors.Is(pathErr.Err, os.ErrNotExist) {
			return err
		}
	}

	err = os.Symlink(origPath, confPath)
	if err != nil {
		return err
	}

	return nil
}

func (s defaultService) applyMonitors(ctx context.Context, config MonitorMap) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var (
		monitors   []hyprland.Monitor          = make([]hyprland.Monitor, 0, len(config)+1)
		workspaces []hyprland.MonitorWorkspace = make([]hyprland.MonitorWorkspace, 0)
	)

	logger.Info("Adding catch all monitor rule")
	monitors = append(monitors, hyprland.Monitor{
		Name:       "",
		Resolution: "preferred",
		Position:   "auto",
		Scale:      "auto",
		Enabled:    true,
	})

	logger.Info("Converting to Hyprland monitor")
	for name, spec := range config {
		var resolution string

		if spec.Resolution != "preferred" && spec.Frequency != "" {
			resolution = fmt.Sprintf("%s@%s", spec.Resolution, spec.Frequency)
		} else {
			resolution = spec.Resolution
		}

		monitors = append(monitors, hyprland.Monitor{
			Name:       name,
			Resolution: resolution,
			Position:   spec.Position,
			Scale:      spec.Scale,
			Enabled:    spec.Enabled,
		})

		if !spec.Enabled {
			continue
		}

		for _, workspaceProfile := range spec.Workspaces {
			workspaces = append(workspaces, hyprland.MonitorWorkspace{
				Monitor:    name,
				ID:         workspaceProfile.ID,
				Default:    workspaceProfile.Default,
				Persistent: workspaceProfile.Persistent,
				Decorate:   workspaceProfile.Decorate,
			})
		}
	}

	logger.Info("Applying Hyprland configuration")
	logger.Debug("Hyprland monitors data", slog.Any("data", monitors))
	logger.Debug("Hyprland workspaces data", slog.Any("data", workspaces))
	return s.hyprland.Apply(ctx, monitors, workspaces)
}

func (s defaultService) applyPanels(ctx context.Context, config Config) error {
	var (
		logger *slog.Logger
		layout hyprpanel.BarLayout
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	logger.Info("Assigning panels to monitors")
	layout, err = assignMonitorPanels(config)
	if err != nil {
		return err
	}

	logger.Info("Applying Hyprpanel configuration")
	logger.Debug("Hyprpanel layout data", slog.Any("data", layout))
	return s.hyprpanel.Apply(ctx, layout)
}

func assignMonitorPanels(config Config) (hyprpanel.BarLayout, error) {
	var layout hyprpanel.BarLayout = make(hyprpanel.BarLayout, len(config.Monitors))

	for _, monitor := range config.Monitors {
		var spec panelSpec

		if monitor.Main {
			spec = config.Panels[keyDefaultPanelMain]
		} else {
			spec = config.Panels[keyDefaultPanelSub]
		}

		layout[monitor.ID] = hyprpanel.BarWidgetConfig{
			L: spec.L,
			R: spec.R,
			M: spec.M,
		}
	}

	return layout, nil
}
