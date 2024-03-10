package database

type Database interface {
	GetRepositories() (*Repositories, error)
	Close() error
}
