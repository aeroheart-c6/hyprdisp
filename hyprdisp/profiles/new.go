package profiles

import (
	"context"

	"aeroheart.io/hyprdisp/hypr"
)

type Controller interface {
	Detect(context.Context, []hypr.Monitor) bool
	Define(context.Context, []hypr.Monitor) error
}

type ControllerImpl struct {
}
