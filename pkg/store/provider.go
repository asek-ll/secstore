package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/asek-ll/secstore/pkg/store/adapter"
)

type config struct {
	Engine    string             `json:"engine"`
	GpgConfig *adapter.GPGConfig `json:"gpgConfig"`
}

func readConfig() (*config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := path.Join(homeDir, ".config/secstore/config.json")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(configPath)
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var conf config

	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func GetDefaultStore() (Store, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}
	if config != nil {
		switch config.Engine {
		case "macos-security":
			return adapter.MacSecurityAdapter{}, nil
		case "gpg":
			store, err := adapter.NewGPGAdapter(config.GpgConfig)
			if err != nil {
				return nil, err
			}

			return store, nil
		}
	}
	if runtime.GOOS == "darwin" {
		return adapter.MacSecurityAdapter{}, nil
	}
	return nil, fmt.Errorf("No security store adapter defined for OS %s", runtime.GOOS)
}
