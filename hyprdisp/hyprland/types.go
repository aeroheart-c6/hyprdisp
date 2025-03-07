package hyprland

import (
	"fmt"
	"strconv"
	"strings"
)

type Event struct {
	Name EventName
	Data []string
}

type MonitorWorkspace struct {
	ID         string
	Monitor    string
	Default    bool
	Persistent bool
	Decorate   bool
}

func (m MonitorWorkspace) marshal() []string {
	return []string{
		fmt.Sprintf("workspace = %s, monitor:%s", m.ID, m.Monitor),
		fmt.Sprintf("workspace = %s, default:%t", m.ID, m.Default),
		fmt.Sprintf("workspace = %s, persistent:%t", m.ID, m.Persistent),
		fmt.Sprintf("workspace = %s, decorate:%t", m.ID, m.Decorate),
	}
}

type Monitor struct {
	ID          string
	Name        string
	Resolution  string
	Position    string
	Scale       string
	Description string
	Make        string
	Model       string
	Serial      string
	Enabled     bool
}

func (m *Monitor) set(field string, value string) error {
	switch field {
	case "id":
		m.ID = value
	case "name":
		m.Name = value
	case "resolution":
		m.Resolution = value
	case "position":
		m.Position = value
	case "scale":
		m.Scale = value
	case "description":
		m.Description = value
	case "make":
		m.Make = value
	case "model":
		m.Model = value
	case "serial":
		m.Serial = value
	case "disabled":
		var (
			disabled bool
			err      error
		)

		disabled, err = strconv.ParseBool(strings.ToLower(value))
		if err != nil {
			return err
		}

		m.Enabled = !disabled
	default:
		return fmt.Errorf("unsupported field: %v with value %v", field, value)
	}

	return nil
}

func (m Monitor) marshal() string {
	return fmt.Sprintf("monitor = %s, %s, %s, %s",
		m.Name,
		m.Resolution,
		m.Position,
		m.Scale,
	)
}

func (m Monitor) String() string {
	return fmt.Sprintf("[%s %s %s]",
		m.ID,
		m.Name,
		m.Description,
	)
}

func (m Monitor) IsZero() bool {
	return m.Name == ""
}
