package memorycache

type cache struct{
	data map[string]interface{}
}

func (c *cache) Set(key string, value interface{}) {
	c.data[key] = value	
}

func (c *cache) Get(key string) interface{} {
	return c.data[key]
}

// New constructs a basic implementation of an in-memory cache.
// No expiry whatsoever, what you set is there forever.
func New() (*cache, error) {
	return &cache{}, nil
}
