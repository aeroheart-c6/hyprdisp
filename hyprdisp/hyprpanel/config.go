package hyprpanel

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func (s defaultService) Apply(layout BarLayout) error {
	var (
		cfgDirPath string
		err        error
	)

	cfgDirPath, err = os.UserConfigDir()
	if err != nil {
		return err
	}

	// unmarshal current configuration JSON file
	var cfg map[string]any

	cfg, err = unmarshalConfig(path.Join(
		cfgDirPath,
		"hyprpanel",
		"config.json",
	))
	if err != nil {
		return err
	}

	for k, v := range cfg {
		fmt.Printf("%v == %+v\n", k, v)
	}

	// add BarLayout instance into the map

	// write the configuration file

	return nil
}

func unmarshalConfig(cfgPath string) (map[string]any, error) {
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
