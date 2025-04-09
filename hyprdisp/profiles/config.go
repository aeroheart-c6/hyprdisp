package profiles

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path"

	"crypto/sha3"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
	"github.com/pelletier/go-toml/v2"
)

// Detect will check if the appropriate configuration file for the active monitors already exists
func (s defaultService) Detect(ctx context.Context, monitors []hyprland.Monitor) (Config, error) {
	var (
		profileID   string = getProfileID(monitors)
		profilePath string
		err         error
	)

	profilePath, err = s.getProfilePath(profileID)
	if err != nil {
		return Config{}, err
	}

	_, err = os.Stat(profilePath)
	if err != nil {
		return Config{}, err
	}

	return s.loadProfile(ctx, profileID)
}

// Init will create a set of config files with default values based on `hyprctl monitors` output
func (s defaultService) Init(ctx context.Context, hyprMonitors []hyprland.Monitor) (Config, error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return Config{}, err
	}

	var (
		monitors monitorConfig = make(monitorConfig, len(hyprMonitors))
		devices  []deviceSpec  = make([]deviceSpec, 0, len(hyprMonitors))
		config   Config
	)
	for _, monitor := range hyprMonitors {
		logger.Info("Found monitor",
			slog.Any("monitor", monitor),
			slog.Bool("enabled", monitor.Enabled),
		)

		monitors[monitor.Name] = monitorSpec{
			ID:         monitor.ID,
			Main:       monitor.ID == "0",
			Position:   "auto",
			Scale:      "auto",
			Resolution: "preferred",
			Frequency:  "",
			Workspaces: []workspaceSpec{
				{
					ID:         fmt.Sprintf("%s001", monitor.ID),
					Default:    true,
					Persistent: true,
					Decorate:   true,
				},
			},
		}

		devices = append(devices, deviceSpec{
			ID:          monitor.ID,
			Name:        monitor.Name,
			Description: monitor.Description,
			Serial:      monitor.Serial,
		})
	}

	config = Config{
		Devices:  devices,
		Monitors: monitors,
		Panels: panelProfile{
			keyDefaultPanelMain: panelSpec{
				L: []string{
					"workspaces",
					"windowtitle",
				},
				M: []string{},
				R: []string{
					"volume",
					"network",
					"bluetooth",
					"systray",
					"clock",
					"dashboard",
				},
			},
			keyDefaultPanelSub: panelSpec{
				L: []string{
					"workspaces",
					"windowtitle",
				},
				M: []string{},
				R: []string{},
			},
		},
	}

	err = s.saveProfile(ctx, getProfileID(hyprMonitors), config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (s defaultService) Apply(ctx context.Context, cfg Config) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	err = s.applyPanels(ctx, cfg)
	if err != nil {
		logger.Info("Unable to apply panel configuration", slog.Any("error", err))
	}

	err = s.applyMonitors(ctx, cfg.Monitors)
	if err != nil {
		return err // TODO should probably try to roll back???
	}

	return nil
}

func (s defaultService) saveProfile(ctx context.Context, id string, profile Config) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var body []byte
	body, err = toml.Marshal(profile)
	if err != nil {
		return err
	}

	var (
		filePath string
		file     *os.File
	)
	filePath, err = s.getProfilePath(id)
	if err != nil {
		return err
	}

	file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logger.Info("Saving TOML data to file", slog.String("file", string(body)))

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (s defaultService) loadProfile(ctx context.Context, id string) (Config, error) {
	var (
		logger   *slog.Logger
		filePath string
		data     []byte
		err      error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return Config{}, err
	}

	filePath, err = s.getProfilePath(id)
	if err != nil {
		return Config{}, err
	}

	data, err = os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var profile Config
	err = toml.Unmarshal(data, &profile)
	if err != nil {
		return Config{}, err
	}

	logger.Info("Loaded profile from TOML", slog.Any("profile", profile))

	return profile, nil
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

func getProfileID(monitors []hyprland.Monitor) string {
	var hash *sha3.SHA3 = sha3.New256()

	for _, monitor := range monitors {
		hash.Write([]byte(monitor.String()))
	}

	return hex.EncodeToString(hash.Sum(nil))[:10]
}
