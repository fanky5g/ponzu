package repositories

type Cacheable interface {
	Cache() Cache
	InvalidateCache() error
}
