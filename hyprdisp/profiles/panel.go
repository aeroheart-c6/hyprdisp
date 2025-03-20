package profiles

import (
	"context"

	"aeroheart.io/hyprdisp/hyprpanel"
)

func (s defaultService) applyPanels(ctx context.Context, profile panelProfile) error {
	var layout hyprpanel.BarLayout = make(hyprpanel.BarLayout, len(profile))

	for key, config := range profile {
		layout[key] = hyprpanel.BarWidgetConfig{
			L: config.L,
			R: config.R,
			M: config.M,
		}
	}

	return s.hyprpanel.Apply(ctx, layout)
}
