package profiles

import (
	"context"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
)

const (
	cfgDirectory string = "hyprdisp"
)

type Service interface {
	Detect(context.Context, []hyprland.Monitor) bool
	Init(context.Context, []hyprland.Monitor) error
	ApplyMonitors(context.Context) error
	ApplyPanels(context.Context) error
	ListenEvents(context.Context, chan error, chan hyprland.Event)
	ListenTimer(context.Context, chan error)
}

type defaultService struct {
	hyprland hyprland.Service
	cfgPath  string
	timer    *time.Timer
	state    state
}

func (s defaultService) ApplyMonitors(context.Context) error {
	return nil
}

func NewDefaultService(hyprlandService hyprland.Service) Service {
	var c defaultService = defaultService{
		hyprland: hyprlandService,
		timer:    time.NewTimer(0),
		state:    "",
	}

	return c
}
