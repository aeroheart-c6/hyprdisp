package hypr

import (
	"errors"
	"os"
	"path"
)

func getInstanceID() (string, error) {
	var (
		value string
		found bool
	)

	value, found = os.LookupEnv(envHyprlandInstanceSignature)
	if !found {
		return "", errors.New("hyprland instance signature not found")
	}

	return value, nil
}

func getRuntimePath() (string, error) {
	var (
		value string
		found bool
	)

	value, found = os.LookupEnv(envXDGRuntimeDirectory)
	if !found {
		return "", errors.New("hyprland runtime path could not be acquired")
	}

	return path.Join(value, "hypr"), nil

}
