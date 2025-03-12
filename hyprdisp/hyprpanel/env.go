package hyprpanel

import (
	"os"
	"path"
)

func (s defaultService) getConfigPath() (string, error) {
	if s.overrideConfigPath != "" {
		return s.overrideConfigPath, nil
	}

	var (
		dir string
		err error
	)

	dir, err = os.UserConfigDir()
	if err != nil {
		return "", nil
	}

	return path.Join(dir, "hyprpanel"), nil
}

func (s defaultService) getConfigFilePath() (string, error) {
	var (
		dir string
		err error
	)

	dir, err = s.getConfigPath()
	if err != nil {
		return "", nil
	}

	return path.Join(dir, "config.json"), nil
}
