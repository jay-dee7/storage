package storage

import "github.com/mailhog/data"

// Storage represents a storage backend
type Storage interface {
	Store(m *data.Message) (string, error)
	List(start, limit int) (*data.Messages, error)
	Search(kind, query string, start, limit int) (*data.Messages, int, error)
	Count() int
	DeleteOne(id string) error
	DeleteAll() error
	Load(id string) (*data.Message, error)
}

// Storage represents a storage backend
type MultiTenantStorage interface {
	Store(m *data.Message, tenant string) (string, error)
	List(start, limit int, tenant string) (*data.Messages, error)
	Search(kind, query string, start, limit int, tenant string) (*data.Messages, int, error)
	Count(tenant string) int
	DeleteOne(id, tenant string) error
	DeleteAll(tenant string) error
	Load(id, tenant string) (*data.Message, error)
}
