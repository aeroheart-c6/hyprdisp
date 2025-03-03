package profiles

import (
	"io"
	"os"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
)

func Test_parseConfig(t *testing.T) {
	var (
		file *os.File
		data []byte
		err  error
	)

	file, err = os.Open("testdata/profiles.toml")
	if err != nil {
		t.Fatalf("failed reading sample configuration file: %v", err)
	}
	defer func() {
		var err error = file.Close()
		if err != nil {
			t.Fatalf("failed closing the file")
		}
	}()

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatalf("failed ot read the file: %v", err)
	}

	var (
		cfg map[string]displayConfig
	)
	err = toml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		t.Fatalf("failed to parse toml file: %v", err)
	}

	t.Logf("%+v", cfg)
}
