package profiles

import (
	"context"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sort"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
)

func (s defaultService) ConnectedMonitors(ctx context.Context) (MonitorMap, error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return nil, err
	}

	var (
		monitors []hyprland.Monitor
		mapping  MonitorMap
	)

	monitors, err = s.hyprland.GetMonitors()
	if err != nil {
		return nil, err
	}

	mapping = make(MonitorMap, len(monitors))
	for _, monitor := range monitors {
		logger.Info("Found monitor",
			slog.String("id", monitor.ID),
			slog.String("name", monitor.Name),
			slog.Bool("enabled", monitor.Enabled),
		)

		mapping[monitor.Name] = monitorSpec{
			ID:          monitor.ID,
			Main:        monitor.ID == "0",
			Name:        monitor.Name,
			Description: monitor.Description,
			Enabled:     monitor.Enabled,
			Position:    "auto",
			Scale:       "auto",
			Resolution:  "preferred",
			Frequency:   "",
			Workspaces: []workspaceSpec{
				{
					ID:         fmt.Sprintf("%s001", monitor.ID),
					Default:    true,
					Persistent: true,
					Decorate:   true,
				},
			},
		}
	}

	return mapping, nil
}

func (s defaultService) SetupDirectories() error {
	var (
		path string
		err  error
	)

	path, err = s.getConfigPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(path, 0o755)
}

func (s defaultService) getConfigPath() (string, error) {
	if s.cfgPath != "" {
		return s.cfgPath, nil
	}

	var (
		dir string
		err error
	)

	dir, err = os.UserConfigDir()
	if err != nil {
		return "", nil
	}

	return path.Join(dir, cfgDirectory), nil
}

func (s defaultService) getProfilePath(id string) (string, error) {
	var (
		confPath string
		err      error
	)

	confPath, err = s.getConfigPath()
	if err != nil {
		return "", err
	}

	return path.Join(confPath, fmt.Sprintf("%v.toml", id)), nil
}

func getProfileID(monitors MonitorMap) (string, error) {
	var (
		hash  *sha3.SHA3 = sha3.New256()
		lines []string   = make([]string, 0, len(monitors))
		err   error
	)

	for _, monitor := range monitors {
		lines = append(lines, monitor.String())
	}
	sort.Stable(sort.StringSlice(lines))

	for _, line := range lines {
		_, err = hash.Write([]byte(line))
		if err != nil {
			return "", nil
		}
	}

	return hex.EncodeToString(hash.Sum(nil))[:10], nil
}
