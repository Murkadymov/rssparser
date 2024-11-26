package cache

type Cache struct {
	data []string
}

func (c *Cache) Get() []string {
	return c.data
}

func (c *Cache) Set(items []string) {
	c.data = append(c.data, items...)
}
