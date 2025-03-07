package hyprland

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

func getEventsSocketPath() (string, error) {
	var (
		runtimeDir string
		instanceID string
		err        error
	)

	runtimeDir, err = getRuntimePath()
	if err != nil {
		return "", nil
	}

	instanceID, err = getInstanceID()
	if err != nil {
		return "", nil
	}

	return path.Join(
		runtimeDir,
		instanceID,
		".socket2.sock",
	), nil
}

func getCommandsSocketPath() (string, error) {
	var (
		runtimeDir string
		instanceID string
		err        error
	)

	runtimeDir, err = getRuntimePath()
	if err != nil {
		return "", nil
	}

	instanceID, err = getInstanceID()
	if err != nil {
		return "", nil
	}

	return path.Join(
		runtimeDir,
		instanceID,
		".socket.sock",
	), nil
}

func getConfigPath() (string, error) {
	var (
		configDir string
		err       error
	)

	configDir, err = os.UserConfigDir()
	if err != nil {
		return "", nil
	}

	return path.Join(
		configDir,
		"hypr",
	), nil
}
