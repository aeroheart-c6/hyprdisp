package profiles

import (
	"context"
	"fmt"

	"aeroheart.io/hyprdisp/hyprland"
)

func (s defaultService) applyMonitors(ctx context.Context, profile monitorProfile) error {
	var (
		monitors   []hyprland.Monitor          = make([]hyprland.Monitor, 0, len(profile))
		workspaces []hyprland.MonitorWorkspace = make([]hyprland.MonitorWorkspace, 0)
	)

	for name, config := range profile {
		var resolution string

		if config.Resolution != "preferred" && config.Frequency != "" {
			resolution = fmt.Sprintf("%s@%s", config.Resolution, config.Frequency)
		} else {
			resolution = config.Resolution
		}

		monitors = append(monitors, hyprland.Monitor{
			Name:       name,
			Resolution: resolution,
			Position:   config.Position,
			Scale:      config.Scale,
		})

		for _, workspaceProfile := range config.Workspaces {
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
