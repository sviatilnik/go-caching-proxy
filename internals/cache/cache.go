package cache

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"net/http"
)

type Cache struct {
	store  Store
	hasher hash.Hash
	ttl    int
}

func NewCache(store Store, TTL int) *Cache {
	return &Cache{
		store:  store,
		hasher: md5.New(),
		ttl:    TTL,
	}
}

func (c *Cache) Get(r *http.Request) (*Response, error) {
	return c.store.Get(r.Context(), c.createKey(r)), nil
}

func (c *Cache) Has(r *http.Request) bool {
	key := c.createKey(r)
	return c.store.Has(r.Context(), key)
}

func (c *Cache) Save(r *http.Request, resp *Response) error {
	key := c.createKey(r)
	return c.store.Save(r.Context(), key, resp, c.ttl)
}

func (c *Cache) createKey(r *http.Request) string {
	c.hasher.Reset()

	c.hasher.Write([]byte(r.Host + r.URL.Path + r.URL.RawQuery))

	h := hex.EncodeToString(c.hasher.Sum(nil))

	return h
}
