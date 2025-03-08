package profiles

import (
	"context"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
)

type Service interface {
	Detect(context.Context, []hyprland.Monitor) bool
	Define(context.Context, []hyprland.Monitor) error
	LoadPanels(context.Context) error
	ListenEvents(context.Context, chan error, chan hyprland.Event)
	ListenTimer(context.Context, chan error)
}

type defaultService struct {
	hyprland hyprland.Service
	timer    *time.Timer
	state    state
}

func NewDefaultService(hyprlandService hyprland.Service) Service {
	var c defaultService = defaultService{
		hyprland: hyprlandService,
		timer:    time.NewTimer(0),
		state:    "",
	}

	return c
}
