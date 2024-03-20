package adapter

import "fmt"

type NoExistsKeyError struct {
	key string
}

func (m *NoExistsKeyError) Error() string {
	return fmt.Sprintf("No key '%s' present in store", m.key)
}
