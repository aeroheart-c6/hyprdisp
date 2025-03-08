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

func (s defaultService) Detect(ctx context.Context, monitors []hyprland.Monitor) bool {
	var (
		configID   string = getDisplaysConfigID(monitors)
		configPath string = getDisplaysConfigPath(ctx, configID)
		err        error
	)

	_, err = os.Stat(configPath)
	return err == nil
}

// Define will create a set of config files with default values based on `hyprctl monitors` output
func (s defaultService) Define(ctx context.Context, monitors []hyprland.Monitor) error {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	var displays displayProfile = make(displayProfile, len(monitors))
	for _, monitor := range monitors {
		logger.Printf("Found monitor (enabled: %v): %+v", monitor.Enabled, monitor)

		displays[monitor.Name] = displayConfig{
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

	return saveDisplaysConfig(ctx, getDisplaysConfigID(monitors), displays)
}

func saveDisplaysConfig(ctx context.Context, id string, displays displayProfile) error {
	var (
		logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		body   []byte
		err    error
	)

	body, err = toml.Marshal(displays)
	if err != nil {
		return err
	}

	var (
		filePath string = getDisplaysConfigPath(ctx, id)
		file     *os.File
	)
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

func loadDisplaysConfig(ctx context.Context, id string) (displayProfile, error) {
	var (
		filePath string = getDisplaysConfigPath(ctx, id)
		data     []byte
		err      error
	)

	data, err = os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var displays displayProfile
	err = toml.Unmarshal(data, &displays)
	if err != nil {
		return nil, err
	}

	return displays, nil
}

func getDisplaysConfigPath(ctx context.Context, id string) string {
	return path.Join(".", "var", fmt.Sprintf("%v.toml", id))
}

func getDisplaysConfigID(monitors []hyprland.Monitor) string {
	var hash *sha3.SHA3 = sha3.New256()

	for _, monitor := range monitors {
		hash.Write([]byte(monitor.String()))
	}

	return hex.EncodeToString(hash.Sum(nil))[:10]
}
