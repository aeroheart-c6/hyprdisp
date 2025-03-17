package profiles

type monitorProfile map[string]monitorConfig

type monitorConfig struct {
	ID         string            `toml:"id"`
	Main       bool              `toml:"main"`
	Resolution string            `toml:"resolution"`
	Frequency  string            `toml:"frequency"`
	Scale      string            `toml:"scale"`
	Workspaces []workspaceConfig `toml:"workspaces"`
}

type workspaceConfig struct {
	ID         string `toml:"id"`
	Default    bool   `toml:"default"`
	Persistent bool   `toml:"persistent"`
	Decorate   bool   `toml:"decorate"`
}

type panelProfile map[string]panelConfig

type panelConfig struct {
	L []string `toml:"left"`
	R []string `toml:"right"`
	M []string `toml:"middle"`
}

type Config struct {
	Panels   panelProfile   `toml:"panels"`
	Monitors monitorProfile `toml:"monitors"`
}

func (p Config) IsZero() bool {
	return p.Panels == nil && p.Monitors == nil
}
