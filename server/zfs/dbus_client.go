package zfs

import (
	"context"
	"encoding/json"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/nickrobison/terraform-linux-provider/server/bus"
	"github.com/nickrobison/terraform-linux-provider/server/middleware"
	"github.com/rs/zerolog"
)

var (
	destination = "com.nickrobison.dbus.zfs1"
	pathname    = "/com/nickrobison/dbus/zfs1"
	prefix      = "com.nickrobison.dbus.ZFS1."
)

type ZfsDebusClient struct {
	conn *dbus.Conn
	log  *zerolog.Logger
	obj  dbus.BusObject
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

	log = log.With().Interface("path", obj.Path()).Logger()

	return &ZfsDebusClient{conn: conn, obj: obj, log: &log}, nil
}

func (c *ZfsDebusClient) ListPools(ctx context.Context) ([]*ZpoolObject, error) {
	m := prefix + "Pools"
	var poolObjs []dbus.ObjectPath
	err := c.obj.CallWithContext(ctx, m, 0).Store(&poolObjs)
	if err != nil {
		return nil, err
	}
	c.log.Debug().Interface("paths", poolObjs).Msg("Received zpools")

	pools := make([]*ZpoolObject, len(poolObjs))
	for i, p := range poolObjs {
		obj := c.conn.Object(destination, p)
		pools[i] = NewZpoolObject(obj, c.log)
	}

	return pools, nil

}

func (c *ZfsDebusClient) Version() (string, error) {
	name := prefix + "Version"
	version, err := bus.Decode[string](c.log, c.obj, name)
	return version, err
}
