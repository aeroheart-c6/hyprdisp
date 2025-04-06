package profiles

import (
	"context"
	"log/slog"

	"aeroheart.io/hyprdisp/hyprpanel"
	"aeroheart.io/hyprdisp/sys"
)

const (
	keyDefaultPanelMain string = "main"
	keyDefaultPanelSub  string = "sub"
)

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

	for name, monitorConfig := range config.Monitors {
		var panelConfig panelConfig

		if monitorConfig.Main {
			panelConfig = config.Panels[keyDefaultPanelMain]
		} else {
			panelConfig = config.Panels[keyDefaultPanelSub]
		}

		layout[name] = hyprpanel.BarWidgetConfig{
			L: panelConfig.L,
			R: panelConfig.R,
			M: panelConfig.M,
		}
	}

	return layout, nil
}
