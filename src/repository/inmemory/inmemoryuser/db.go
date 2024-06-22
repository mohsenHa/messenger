package inmemoryuser

import (
	"sync"
)

type DB struct {
	table sync.Map
}

func New() *DB {
	return &DB{
		table: sync.Map{},
	}
}
