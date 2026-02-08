package firewall

import (
	"context"
)

type FirewallClient interface {
	ListZones(ctx context.Context) ([]*ZoneObject, error)
	GetZone(ctx context.Context, name string) (*ZoneObject, error)
	AddZone(ctx context.Context, name string, settings ZoneSettings) error
	RemoveZone(ctx context.Context, name string) error
	Version() (string, error)
}

type ZoneSettings struct {
	Description string
	Target      string
}
