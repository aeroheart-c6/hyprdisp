package profiles

import (
	"context"

	"aeroheart.io/hyprdisp/hyprland"
)

type Controller interface {
	Detect(context.Context, []hyprland.Monitor) bool
	Define(context.Context, []hyprland.Monitor) error
	LoadPanels(context.Context) error
}

type ControllerImpl struct {
}
