package hyprpanel

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"aeroheart.io/hyprdisp/sys"
)

func (s defaultService) Apply(
	ctx context.Context,
	layout BarLayout,
) error {
	var (
		logger  *slog.Logger
		cfgPath string
		cfg     map[string]any
		err     error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	cfgPath, err = s.getConfigFilePath()
	if err != nil {
		return err
	}

	logger.Info("Loading configuration file", slog.String("path", cfgPath))
	cfg, err = loadConfig(cfgPath)
	if err != nil {
		return err
	}

	// add BarLayout instance into the map
	cfg[cfgKeyLayouts] = layout

	logger.Info("Writing configuration file", slog.String("path", cfgPath))
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
	if errors.Is(err, os.ErrNotExist) {
		return map[string]any{}, nil
	} else if err != nil {
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
