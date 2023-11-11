package middleware

import "context"

type Cacher interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (string, error)
}

type Storager interface {
	GetAlias(ctx context.Context, orig string) (string, error)
	GetOriginal(ctx context.Context, alias string) (string, error)
	Set(ctx context.Context, alias string, orig string) error
}
