package profiles

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sort"
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

func (s defaultService) getProfilePath(id string) (string, error) {
	var (
		confPath string
		err      error
	)

	confPath, err = s.getConfigPath()
	if err != nil {
		return "", err
	}

	return path.Join(confPath, fmt.Sprintf("%v.toml", id)), nil
}

func getProfileID(monitors MonitorMap) (string, error) {
	var (
		hash  *sha3.SHA3 = sha3.New256()
		lines []string   = make([]string, 0, len(monitors))
		err   error
	)

	for _, monitor := range monitors {
		lines = append(lines, monitor.String())
	}
	sort.Stable(sort.StringSlice(lines))

	for _, line := range lines {
		_, err = hash.Write([]byte(line))
		if err != nil {
			return "", nil
		}
	}

	return hex.EncodeToString(hash.Sum(nil))[:10], nil
}
