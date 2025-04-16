package profiles

import (
	"context"
	"fmt"

	"aeroheart.io/hyprdisp/hyprland"
)

func (s defaultService) applyMonitors(ctx context.Context, config monitorConfig) error {
	var (
		monitors   []hyprland.Monitor          = make([]hyprland.Monitor, 0, len(config))
		workspaces []hyprland.MonitorWorkspace = make([]hyprland.MonitorWorkspace, 0)
	)

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

	return s.hyprland.Apply(ctx, monitors, workspaces)
}
