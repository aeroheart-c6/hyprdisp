package hyprland

import (
	"context"
	"path"
)

type Service interface {
	GetMonitors() ([]Monitor, error)
	Apply([]Monitor, []MonitorWorkspace) error
	StreamEvents(context.Context) (chan Event, chan error, error)
}

type defaultService struct {
	overrideConfigPath string
}

func NewDefaultService() Service {
	return defaultService{
		overrideConfigPath: path.Join(".", "var"),
	}
}
