package cache

type Cache struct {
	data []string
}

func (c *Cache) Get() []string {}

func (c *Cache) Set(items []string) {}
