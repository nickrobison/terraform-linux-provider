package zfs

import (
	"context"

	"github.com/nickrobison/terraform-linux-provider/common"
)

type ZfsClient interface {
	ListPools(ctx context.Context) ([]common.ZPool, error)
	Version() (string, error)
}
