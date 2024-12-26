package driver

type Database interface {
	GetRepositoryByToken(token string) Repository
	Close() error
}
