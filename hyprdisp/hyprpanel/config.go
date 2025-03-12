package hyprpanel

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"aeroheart.io/hyprdisp/sys"
)

func (s defaultService) Apply(
	ctx context.Context,
	layout BarLayout,
) error {
	var (
		logger  *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		cfgPath string
		cfg     map[string]any
		err     error
	)

	cfgPath, err = s.getConfigFilePath()
	if err != nil {
		return err
	}

	logger.Printf("Loading configuration file at: %s", cfgPath)
	cfg, err = loadConfig(cfgPath)
	if err != nil {
		return err
	}

	// add BarLayout instance into the map
	cfg[cfgKeyLayouts] = layout

	logger.Printf("Writing configuration file at: %s", cfgPath)
	err = writeConfig(cfgPath, cfg)
	if err != nil {
		return err
	}

	return nil
}

func loadConfig(cfgPath string) (map[string]any, error) {
	var (
		data []byte
		err  error
	)

	data, err = os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	var cfg map[string]any
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func writeConfig(cfgPath string, cfg map[string]any) error {
	var (
		file *os.File
		data []byte
		err  error
	)
	file, err = os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	data, err = json.MarshalIndent(cfg, "", cfgIndent)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
