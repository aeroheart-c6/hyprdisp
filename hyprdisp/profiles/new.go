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
	Detect(context.Context, MonitorMap) (Config, error)
	Init(context.Context, MonitorMap) (Config, error)
	Apply(context.Context, Config) error
	ConnectedMonitors(context.Context) (MonitorMap, error)

	SetupDirectories() error
	AsListener() ListenerService
}

type ListenerService interface {
	ListenHyprland(context.Context, chan hyprland.Event, chan error)
	ListenTimer(context.Context, chan error)
}

type defaultService struct {
	hyprland  hyprland.Service
	hyprpanel hyprpanel.Service
	cfgPath   string
	timer     *time.Timer
	state     state
}

func (s defaultService) AsListener() ListenerService {
	return &s
}

func NewDefaultService(
	hyprlandSrv hyprland.Service,
	hyprpanelSrv hyprpanel.Service,
	cfgPath string,
) Service {
	var service defaultService = defaultService{
		hyprland:  hyprlandSrv,
		hyprpanel: hyprpanelSrv,
		timer:     nil,
		state:     WatchState,
		cfgPath:   cfgPath,
	}

	return service
}
