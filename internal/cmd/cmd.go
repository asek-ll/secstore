package cmd

type Store interface {
	Load(key string) (string, error)
	Save(key string, value string) error
	Delete(key string) error
}

type Config struct {
	Store Store
}

func (cfg *Config) SetConfig(source *Config) {
	cfg.Store = source.Store
}

type Command interface {
	SetConfig(source *Config)
	Execute(args []string) error
}
