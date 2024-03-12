package store

type Store interface {
	Load(key string) (string, error)
	Save(key string, value string) error
	Delete(key string) error
}
