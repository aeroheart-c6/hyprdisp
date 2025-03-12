package hyprpanel

type Service interface {
	Apply(BarLayout) error
}

type defaultService struct {
	overrideConfigPath string
}

func NewDefaultService() Service {
	return defaultService{
		overrideConfigPath: "var",
	}
}
