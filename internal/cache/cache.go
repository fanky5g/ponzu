package cache

type Cache interface {
	Set(string, interface{})
	Get(string) interface{}
}
