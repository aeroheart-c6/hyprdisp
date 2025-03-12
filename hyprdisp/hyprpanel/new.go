package hyprpanel

import "context"

const (
	cfgKeyLayouts string = "bar.layouts"
	cfgIndent     string = "    "
	cfgDirectory  string = "hyprpanel"
)

type Service interface {
	Apply(context.Context, BarLayout) error
}

type defaultService struct {
	cfgPath string
	cfgFile string
}

func NewDefaultService() Service {
	return defaultService{}
}
