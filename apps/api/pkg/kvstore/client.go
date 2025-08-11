package kvstore

import "context"

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, opts ...SetOptions) error
	Delete(ctx context.Context, key string) error
}

type SetOptions struct {
	Expiration int64
}
