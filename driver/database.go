package driver

import "github.com/fanky5g/ponzu/tokens"

type Database interface {
	Get(token tokens.Repository) interface{}
	Close() error
}
