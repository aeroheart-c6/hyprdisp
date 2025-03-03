package hyprpanel

type BarLayout map[string]BarWidgetConfig

func (b *BarLayout) Set(monitorID string, widgets BarWidgetConfig) {
	(*b)[monitorID] = widgets
}

type BarWidgetConfig struct {
	L []string `json:"left"`
	R []string `json:"right"`
	M []string `json:"middle"`
}
