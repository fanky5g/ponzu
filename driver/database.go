package driver

type Database interface {
	Get(token string) interface{}
	Close() error
}
