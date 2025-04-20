package profiles

import (
	"context"
	"fmt"
	"log/slog"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/hyprpanel"
	"aeroheart.io/hyprdisp/sys"
)

const (
	keyDefaultPanelMain string = "main"
	keyDefaultPanelSub  string = "sub"
)

func (s defaultService) Apply(ctx context.Context, cfg Config) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	err = s.applyPanels(ctx, cfg)
	if err != nil {
		logger.Info("Unable to apply panel configuration", slog.Any("error", err))
	}

	err = s.applyMonitors(ctx, cfg.Monitors)
	if err != nil {
		return err // TODO should probably try to roll back???
	}

	return nil
}

func (s defaultService) applyMonitors(ctx context.Context, config monitorConfig) error {
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

	logger.Info("Converting to hyprland monitor")
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

	layout, err = assignMonitorPanels(config)
	if err != nil {
		return err
	}

	logger.Info("Applying Hyprpanel configuration", slog.Any("layout", layout))
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
