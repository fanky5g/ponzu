package database

type Database interface {
	GetRepository(name string) Repository
	Close() error
}
