package profiles

type monitorConfig map[string]monitorSpec

type monitorSpec struct {
	ID         string          `toml:"id"`
	Main       bool            `toml:"main"`
	Enabled    bool            `toml:"enabled"`
	Resolution string          `toml:"resolution"`
	Frequency  string          `toml:"frequency"`
	Scale      string          `toml:"scale"`
	Position   string          `toml:"position"`
	Workspaces []workspaceSpec `toml:"workspaces"`
}

type workspaceSpec struct {
	ID         string `toml:"id"`
	Default    bool   `toml:"default"`
	Persistent bool   `toml:"persistent"`
	Decorate   bool   `toml:"decorate"`
}

type panelProfile map[string]panelSpec

type panelSpec struct {
	L []string `toml:"left"`
	R []string `toml:"right"`
	M []string `toml:"middle"`
}

type deviceSpec struct {
	ID          string `toml:"id"`
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Serial      string `toml:"serial"`
}

type Config struct {
	Devices  []deviceSpec  `toml:"devices"`
	Panels   panelProfile  `toml:"panels"`
	Monitors monitorConfig `toml:"monitors"`
}

func (p Config) IsZero() bool {
	return p.Panels == nil && p.Monitors == nil
}
