package hyprland

import (
	"errors"
	"os"
	"path"
)

const (
	envHyprlandInstanceSignature = "HYPRLAND_INSTANCE_SIGNATURE"
	envXDGRuntimeDirectory       = "XDG_RUNTIME_DIR"
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

	return path.Join(value, cfgDirectory), nil

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

	return path.Join(runtimeDir, instanceID, ".socket2.sock"), nil
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

	return path.Join(runtimeDir, instanceID, ".socket.sock"), nil
}

/*
 * TODO make testable
 *
 * it is currently a pain to test because the working directory of the test is the directory of the test file. This
 * means relative file operations will complain because on normal mode, it can find `./var` but not `./testdata` -- BUT
 * testing mode is the reverse. It can find `./testdata` but not `./var`
 *
 * make the program more dynamic in this "configuration" value
 */
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
