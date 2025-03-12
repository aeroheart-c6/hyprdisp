package hyprpanel

import (
	"encoding/json"
	"os"
)

func (s defaultService) Apply(layout BarLayout) error {
	var (
		cfgPath string
		cfg     map[string]any
		err     error
	)

	cfgPath, err = s.getConfigFilePath()
	if err != nil {
		return err
	}

	// unmarshal current configuration JSON file
	cfg, err = loadConfig(cfgPath)
	if err != nil {
		return err
	}

	// add BarLayout instance into the map
	cfg["bar.layouts"] = layout

	// write the configuration file
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

	data, err = json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
