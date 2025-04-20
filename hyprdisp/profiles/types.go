package profiles

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type monitorSpec struct {
	ID          string          `toml:"id"`
	Main        bool            `toml:"main"`
	Name        string          `toml:"name"`
	Description string          `toml:"description"`
	Enabled     bool            `toml:"enabled"`
	Resolution  string          `toml:"resolution"`
	Frequency   string          `toml:"frequency"`
	Scale       string          `toml:"scale"`
	Position    string          `toml:"position"`
	Workspaces  []workspaceSpec `toml:"workspaces"`
}

func (m monitorSpec) String() string {
	return fmt.Sprintf("[%s %s %s]",
		m.ID,
		m.Name,
		m.Description,
	)
}

type workspaceSpec struct {
	ID         string `toml:"id"`
	Default    bool   `toml:"default"`
	Persistent bool   `toml:"persistent"`
	Decorate   bool   `toml:"decorate"`
}

type panelSpec struct {
	L []string `toml:"left"`
	R []string `toml:"right"`
	M []string `toml:"middle"`
}

type MonitorMap map[string]monitorSpec
type PanelMap map[string]panelSpec

type Config struct {
	ID       string     `toml:"-"`
	Panels   PanelMap   `toml:"panels"`
	Monitors MonitorMap `toml:"monitors"`
}

func (c Config) IsZero() bool {
	return c.Panels == nil && c.Monitors == nil
}

func (c Config) ToTOML() ([]byte, error) {
	var (
		headers []string = make([]string, 0, len(c.Monitors))
		body    []byte
		err     error
	)

	for _, device := range c.Monitors {
		headers = append(
			headers,
			fmt.Sprintf("#    * %s", device.String()),
		)
	}
	sort.Stable(sort.StringSlice(headers))

	headers = append(
		[]string{"# Quick summary of monitors in this configuration:"},
		headers...,
	)

	body, err = toml.Marshal(c)
	if err != nil {
		return nil, err
	}

	return append(
		[]byte(strings.Join(headers, "\n")),
		body...,
	), nil
}
