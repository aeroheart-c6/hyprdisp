package hyprpanel

type Service interface {
	Apply(BarLayout) error
}

type defaultService struct {
}

func NewDefaultService() Service {
	return defaultService{}
}
