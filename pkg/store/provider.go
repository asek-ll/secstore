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

type provider struct {
	config       *config
	configLoaded bool
}

func NewProvider() provider {
	return provider{
		config:       nil,
		configLoaded: false,
	}
}

func (ctx provider) readConfig() (*config, error) {
	if ctx.configLoaded {
		return ctx.config, nil
	}

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

	ctx.config = &conf
	ctx.configLoaded = true

	return &conf, nil
}

func (p provider) GetStoreWithEngine(engine string) (Store, error) {
	switch engine {
	case "macos-security":
		return adapter.MacSecurityAdapter{}, nil
	case "gpg":
		config, err := p.readConfig()
		if err != nil {
			return nil, err
		}
		if config == nil {
			return nil, fmt.Errorf("No gpg config present")
		}
		store, err := adapter.NewGPGAdapter(config.GpgConfig)
		if err != nil {
			return nil, err
		}

		return store, nil
	}

	return nil, fmt.Errorf("No security store adapter defined for engine %s", engine)
}

func (p provider) GetDefaultStore() (Store, error) {
	config, err := p.readConfig()
	if err != nil {
		return nil, err
	}
	if config != nil && config.Engine != "" {
		return p.GetStoreWithEngine(config.Engine)
	}
	if runtime.GOOS == "darwin" {
		return adapter.MacSecurityAdapter{}, nil
	}
	return nil, fmt.Errorf("No security store adapter defined for OS %s", runtime.GOOS)
}
