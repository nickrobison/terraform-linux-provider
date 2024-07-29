package zfs

import (
	"context"
	"encoding/json"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/nickrobison/terraform-linux-provider/common"
	"github.com/nickrobison/terraform-linux-provider/server/middleware"
	"github.com/rs/zerolog"
)

var (
	destination = "com.nickrobison.dbus.zfs1"
	pathname    = "/com/nickrobison/dbus/zfs1"
	prefix      = "com.nickrobison.dbus.ZFS1."
)

type ZfsDebusClient struct {
	log *zerolog.Logger
	obj dbus.BusObject
}

func NewZfsClient(conn *dbus.Conn) (ZfsClient, error) {
	log := middleware.Logger()
	log.Info().Msg("Initializing ZFS DBus connection")
	obj := conn.Object(destination, dbus.ObjectPath(pathname))
	node, err := introspect.Call(obj)
	if err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(node, "", "    ")
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("%s\n", string(data))

	return &ZfsDebusClient{obj: obj, log: &log}, nil
}

func (c *ZfsDebusClient) ListPools(ctx context.Context) ([]common.ZPool, error) {
	return nil, nil
}

func (c *ZfsDebusClient) Version() (string, error) {
	var version string
	name := prefix + "Version"
	prop, err := c.obj.GetProperty(name)
	if err != nil {
		return name, err
	}
	c.log.Info().Str("property", name).Interface("variant", prop.Signature().String()).Msg("Received response")
	prop.Store(&version)
	return version, nil
}
