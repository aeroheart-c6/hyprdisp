package profiles

import (
	"context"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/hyprpanel"
)

const (
	cfgDirectory string = "hyprdisp"
)

type Service interface {
	Detect(context.Context, []hyprland.Monitor) (Config, error)
	Init(context.Context, []hyprland.Monitor) (Config, error)
	Apply(context.Context, Config) error
	ListenEvents(context.Context, chan error, chan hyprland.Event)
	ListenTimer(context.Context, chan error)
}

type defaultService struct {
	hyprland  hyprland.Service
	hyprpanel hyprpanel.Service
	cfgPath   string
	timer     *time.Timer
	state     state
}

func NewDefaultService(
	hyprlandSrv hyprland.Service,
	hyprpanelSrv hyprpanel.Service,
	cfgPath string,
) Service {
	var c defaultService = defaultService{
		hyprland:  hyprlandSrv,
		hyprpanel: hyprpanelSrv,
		timer:     time.NewTimer(0),
		state:     "",
		cfgPath:   cfgPath,
	}

	return c
}
