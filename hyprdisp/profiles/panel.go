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
