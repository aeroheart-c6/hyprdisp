package profiles

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"crypto/sha3"

	"aeroheart.io/hyprdisp/hypr"
	"aeroheart.io/hyprdisp/sys"
	"github.com/pelletier/go-toml/v2"
)

type displayProfile map[string]displayConfig

type displayConfig struct {
	ID         string            `toml:"id"`
	Main       bool              `toml:"main"`
	Scale      string            `toml:"scale"`
	Resolution string            `toml:"resolution"`
	Frequency  string            `toml:"frequency"`
	Workspaces []workspaceConfig `toml:"workspaces"`
}

type workspaceConfig struct {
	ID         string `toml:"id"`
	Default    bool   `toml:"default"`
	Persistent bool   `toml:"persistent"`
	Decorate   bool   `toml:"decorate"`
}

func (c ControllerImpl) Detect(ctx context.Context, monitors []hypr.Monitor) bool {
	var (
		idName string = createID(monitors)
		idPath string = path.Join(".", "var", fmt.Sprintf("%v.toml", idName))
		err    error
	)

	_, err = os.Stat(idPath)
	return err == nil
}

// Define will create a set of config files with default values based on `hyprctl monitors` output
func (c ControllerImpl) Define(ctx context.Context, monitors []hypr.Monitor) error {
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

	return writeConfig(ctx, createID(monitors), displays)
}

func writeConfig(ctx context.Context, id string, displays displayProfile) error {
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
		filePath string = path.Join(".", "var", fmt.Sprintf("%v.toml", id))
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

func createID(monitors []hypr.Monitor) string {
	var hash *sha3.SHA3 = sha3.New256()

	for _, monitor := range monitors {
		hash.Write([]byte(monitor.String()))
	}

	return hex.EncodeToString(hash.Sum(nil))[:10]
}
