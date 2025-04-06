package hyprpanel

import (
	"context"
)

const (
	cfgKeyLayouts string = "bar.layouts"
	cfgIndent     string = "    "
	cfgDirectory  string = "hyprpanel"
	cfgFile       string = "config.json"
)

type Service interface {
	Apply(context.Context, BarLayout) error
}

type defaultService struct {
	cfgPath string
	cfgFile string
}

func NewDefaultService(cfgPath string) Service {
	return defaultService{
		cfgPath: cfgPath,
		cfgFile: cfgFile,
	}
}
