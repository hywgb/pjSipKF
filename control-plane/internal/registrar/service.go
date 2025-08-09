package registrar

import (
	"context"
)

type Service interface {
	Register(ctx context.Context, user string, contact string, expiresSeconds int) error
	Deregister(ctx context.Context, user string, contact string) error
	Lookup(ctx context.Context, user string) ([]Binding, error)
}

type Binding struct {
	Contact      string
	ExpiresUnix  int64
}