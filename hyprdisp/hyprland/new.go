package hyprland

import (
	"context"
	"path"
)

const (
	cfgDirectory      string = "hypr"
	cfgMonitorsFile   string = "hypr-monitors.conf"
	cfgWorkspacesFile string = "hypr-workspaces.conf"
)

type Service interface {
	GetMonitors() ([]Monitor, error)
	Apply(context.Context, []Monitor, []MonitorWorkspace) error
	StreamEvents(context.Context) (chan Event, chan error, error)
}

type defaultService struct {
	cfgPath       string
	cfgMonitors   string
	cfgWorkspaces string
}

func NewDefaultService() Service {
	return defaultService{
		cfgPath:       path.Join(".", "var"),
		cfgMonitors:   cfgMonitorsFile,
		cfgWorkspaces: cfgWorkspacesFile,
	}
}
