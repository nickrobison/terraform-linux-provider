package zfs

import (
	"context"
)

type ZfsClient interface {
	ListPools(ctx context.Context) ([]*ZpoolObject, error)
	Version() (string, error)
}
