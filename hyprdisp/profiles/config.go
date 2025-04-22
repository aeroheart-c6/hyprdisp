package profiles

import (
	"context"
	"log/slog"
	"os"

	"aeroheart.io/hyprdisp/sys"
	"github.com/pelletier/go-toml/v2"
)

// Detect will check if the appropriate configuration file for the active monitors already exists
func (s defaultService) Detect(ctx context.Context, monitors MonitorMap) (Config, error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return Config{}, err
	}

	var (
		id   string
		path string
	)
	if monitors == nil {
		monitors, err = s.ConnectedMonitors(ctx)
		if err != nil {
			return Config{}, err
		}
	}

	id, err = getProfileID(monitors)
	if err != nil {
		return Config{}, err
	}

	path, err = s.getProfilePath(id)
	if err != nil {
		return Config{}, err
	}

	logger.Info("Checking if configuration profile exists",
		slog.String("id", id),
		slog.String("path", path),
	)
	_, err = os.Stat(path)
	if err != nil {
		return Config{}, err
	}

	return s.loadProfile(ctx, id)
}

// Init will create a set of config files with default values based on `hyprctl monitors` output
func (s defaultService) Init(ctx context.Context, monitors MonitorMap) (Config, error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if monitors == nil {
		monitors, err = s.ConnectedMonitors(ctx)
		if err != nil {
			return Config{}, err
		}
	}

	config = Config{
		Monitors: monitors,
		Panels: PanelMap{
			keyDefaultPanelMain: panelSpec{
				L: []string{
					"workspaces",
					"windowtitle",
				},
				M: []string{
					"media",
					"notifications",
				},
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

	config.ID, err = getProfileID(config.Monitors)
	if err != nil {
		return Config{}, err
	}

	logger.Info("Initializing profile", slog.String("id", config.ID))
	err = s.saveProfile(ctx, config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// savePRofile writes the Config struct into the configuration file
func (s defaultService) saveProfile(ctx context.Context, config Config) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var (
		data []byte
		path string
		file *os.File
	)

	path, err = s.getProfilePath(config.ID)
	if err != nil {
		return err
	}

	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	data, err = config.toTOML()
	if err != nil {
		return err
	}

	logger.Info("Saving TOML data to file", slog.String("file", path))
	logger.Debug("TOML file data", slog.String("data", string(data)))
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// loadProfile reads the configuration file and unmarshals into the internal Config struct
func (s defaultService) loadProfile(ctx context.Context, id string) (Config, error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return Config{}, err
	}

	var (
		filePath string
		data     []byte
		config   Config
	)
	filePath, err = s.getProfilePath(id)
	if err != nil {
		return Config{}, err
	}

	data, err = os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	config.ID = id

	logger.Info("Loaded profile config from TOML", slog.String("id", config.ID))
	logger.Debug("Config data", slog.Any("data", config))
	return config, nil
}
