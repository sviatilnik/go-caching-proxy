package cache

import "context"

type Store interface {
	Get(ctx context.Context, key string) *Response
	Save(ctx context.Context, key string, value *Response, ttl int) error
	Has(ctx context.Context, key string) bool
	Remove(ctx context.Context, key string) error
}
