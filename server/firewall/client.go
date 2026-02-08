package firewall

import (
	"context"
)

type FirewallClient interface {
	ListZones(ctx context.Context) ([]*ZoneObject, error)
	GetZone(ctx context.Context, name string) (*ZoneObject, error)
	AddZone(ctx context.Context, name string, settings ZoneSettings) error
	RemoveZone(ctx context.Context, name string) error
	AddRichRule(ctx context.Context, zone string, rule string) error
	RemoveRichRule(ctx context.Context, zone string, rule string) error
	AddPort(ctx context.Context, zone string, port string, protocol string) error
	RemovePort(ctx context.Context, zone string, port string, protocol string) error
	AddService(ctx context.Context, zone string, service string) error
	RemoveService(ctx context.Context, zone string, service string) error
	Version() (string, error)
}

type ZoneSettings struct {
	Description string
	Target      string
}
