package hyprpanel

import (
	"os"
	"path"
)

func (s defaultService) getConfigPath() (string, error) {
	if s.cfgPath != "" {
		return s.cfgPath, nil
	}

	var (
		dir string
		err error
	)

	dir, err = os.UserConfigDir()
	if err != nil {
		return "", nil
	}

	return path.Join(dir, cfgDirectory), nil
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

	return path.Join(dir, s.cfgFile), nil
}
