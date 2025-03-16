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
func (s defaultService) Detect(ctx context.Context, monitors []hyprland.Monitor) bool {
	var (
		profileID   string = getProfileID(monitors)
		profilePath string
		err         error
	)

	profilePath, err = s.getProfilePath(profileID)
	if err != nil {
		return false
	}

	_, err = os.Stat(profilePath)
	return err == nil
}

// Init will create a set of config files with default values based on `hyprctl monitors` output
func (s defaultService) Init(ctx context.Context, hyprMonitors []hyprland.Monitor) error {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	var monitors monitorProfile = make(monitorProfile, len(hyprMonitors))
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

	var profile profileConfig = profileConfig{
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

	return s.saveProfile(ctx, getProfileID(hyprMonitors), profile)
}

func (s defaultService) saveProfile(ctx context.Context, id string, profile profileConfig) error {
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

	logger.Printf("Resulting TOML file...")
	logger.Printf("%+v", string(body))

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (s defaultService) loadProfile(ctx context.Context, id string) (monitorProfile, error) {
	var (
		filePath string
		data     []byte
		err      error
	)

	filePath, err = s.getProfilePath(id)
	if err != nil {
		return nil, err
	}

	data, err = os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var displays monitorProfile
	err = toml.Unmarshal(data, &displays)
	if err != nil {
		return nil, err
	}

	return displays, nil
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
