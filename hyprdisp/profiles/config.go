package profiles

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"crypto/sha3"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
	"github.com/pelletier/go-toml/v2"
)

const (
	keyDefaultPanelMain string = "main"
	keyDefaultPanelSub  string = "sub"
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
		logger   *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		monitors monitorProfile
		config   Config
		err      error
	)

	monitors = make(monitorProfile, len(hyprMonitors))
	for _, monitor := range hyprMonitors {
		logger.Printf("Found monitor (enabled: %v): %+v", monitor.Enabled, monitor)

		monitors[monitor.Name] = monitorConfig{
			ID:         monitor.ID,
			Main:       monitor.ID == "0",
			Scale:      "auto",
			Resolution: "preferred",
			Frequency:  "",
			Workspaces: []workspaceConfig{
				{
					ID:         fmt.Sprintf("%s001", monitor.ID),
					Default:    true,
					Persistent: true,
					Decorate:   true,
				},
			},
		}
	}

	config = Config{
		Panels: panelProfile{
			keyDefaultPanelMain: panelConfig{
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
			keyDefaultPanelSub: panelConfig{
				L: []string{
					"workspaces",
					"windowtitle",
				},
				M: []string{},
				R: []string{},
			},
		},
		Monitors: monitors,
	}

	err = s.saveProfile(ctx, getProfileID(hyprMonitors), config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (s defaultService) Apply(ctx context.Context, cfg Config) error {
	var err error

	err = s.applyPanels(ctx, cfg.Panels)
	if err != nil {
		return nil
	}

	err = s.applyMonitors(ctx)
	if err != nil {
		return nil // TODO should probably try to roll back???
	}

	return nil
}

func (s defaultService) saveProfile(ctx context.Context, id string, profile Config) error {
	var (
		logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		body   []byte
		err    error
	)

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

	logger.Printf("Saving TOML data to file: %+v", string(body))

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (s defaultService) loadProfile(ctx context.Context, id string) (Config, error) {
	var (
		logger   *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		filePath string
		data     []byte
		err      error
	)

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

	logger.Printf("Loaded profile from TOML: %+v", profile)

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
