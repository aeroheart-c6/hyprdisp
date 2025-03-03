package profiles

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

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

func (c ControllerImpl) Detect(ctx context.Context, monitors []hypr.Monitor) error {

	return nil
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

	return writeConfig(ctx, displays)
}

func writeConfig(ctx context.Context, displays displayProfile) error {
	var (
		logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		body   []byte
		err    error
	)

	body, err = toml.Marshal(displays)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.OpenFile("config.toml", os.O_CREATE|os.O_WRONLY, 0644)
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

func createID(displays []hypr.Monitor) {
	var identifiers []string = make([]string, 0, len(displays))

	for _, display := range displays {
		identifiers = append(identifiers, display.String())
	}

	var body string = strings.Join(identifiers, "|")
	var hash [32]byte = sha3.Sum256([]byte(body))

	fmt.Printf("Hash of the display combination: %v is %v\n",
		body,
		hex.EncodeToString(hash[:])[:6],
	)
}
