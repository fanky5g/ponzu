package database

type Database interface {
	GetRepositoryByToken(token string) Repository
	Close() error
}
