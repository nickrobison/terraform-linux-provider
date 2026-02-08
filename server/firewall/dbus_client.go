package firewall

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
	destination = "org.fedoraproject.FirewallD1"
	pathname    = "/org/fedoraproject/FirewallD1"
	prefix      = "org.fedoraproject.FirewallD1."
)

type FirewallDebusClient struct {
	conn *dbus.Conn
	log  *zerolog.Logger
	obj  dbus.BusObject
}

func NewFirewallClient(conn *dbus.Conn) (FirewallClient, error) {
	log := middleware.Logger()
	log.Info().Msg("Initializing Firewall DBus connection")
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

	return &FirewallDebusClient{conn: conn, obj: obj, log: &log}, nil
}

func (c *FirewallDebusClient) ListZones(ctx context.Context) ([]*ZoneObject, error) {
	m := prefix + "getZones"
	var zoneNames []string
	err := c.obj.CallWithContext(ctx, m, 0).Store(&zoneNames)
	if err != nil {
		return nil, err
	}
	c.log.Debug().Interface("zones", zoneNames).Msg("Received firewall zones")

	zones := make([]*ZoneObject, len(zoneNames))
	for i, zoneName := range zoneNames {
		// Get zone object path
		m := prefix + "config.getZoneByName"
		var zonePath dbus.ObjectPath
		err := c.obj.CallWithContext(ctx, m, 0, zoneName).Store(&zonePath)
		if err != nil {
			return nil, err
		}
		obj := c.conn.Object(destination, zonePath)
		zones[i] = NewZoneObject(obj, c.log)
	}

	return zones, nil
}

func (c *FirewallDebusClient) GetZone(ctx context.Context, name string) (*ZoneObject, error) {
	m := prefix + "config.getZoneByName"
	var zonePath dbus.ObjectPath
	err := c.obj.CallWithContext(ctx, m, 0, name).Store(&zonePath)
	if err != nil {
		return nil, err
	}
	obj := c.conn.Object(destination, zonePath)
	return NewZoneObject(obj, c.log), nil
}

func (c *FirewallDebusClient) AddZone(ctx context.Context, name string, settings ZoneSettings) error {
	m := prefix + "config.addZone"
	settingsMap := make(map[string]dbus.Variant)
	if settings.Description != "" {
		settingsMap["description"] = dbus.MakeVariant(settings.Description)
	}
	if settings.Target != "" {
		settingsMap["target"] = dbus.MakeVariant(settings.Target)
	}
	return c.obj.CallWithContext(ctx, m, 0, name, settingsMap).Err
}

func (c *FirewallDebusClient) RemoveZone(ctx context.Context, name string) error {
	m := prefix + "config.removeZone"
	return c.obj.CallWithContext(ctx, m, 0, name).Err
}

func (c *FirewallDebusClient) Version() (string, error) {
	name := prefix + "version"
	version, err := bus.Decode[string](c.log, c.obj, name)
	return version, err
}

func (c *FirewallDebusClient) AddRichRule(ctx context.Context, zone string, rule string) error {
	m := prefix + "zone.addRichRule"
	return c.obj.CallWithContext(ctx, m, 0, zone, rule).Err
}

func (c *FirewallDebusClient) RemoveRichRule(ctx context.Context, zone string, rule string) error {
	m := prefix + "zone.removeRichRule"
	return c.obj.CallWithContext(ctx, m, 0, zone, rule).Err
}

func (c *FirewallDebusClient) AddPort(ctx context.Context, zone string, port string, protocol string) error {
	m := prefix + "zone.addPort"
	portProto := port + "/" + protocol
	return c.obj.CallWithContext(ctx, m, 0, zone, portProto, 0).Err
}

func (c *FirewallDebusClient) RemovePort(ctx context.Context, zone string, port string, protocol string) error {
	m := prefix + "zone.removePort"
	portProto := port + "/" + protocol
	return c.obj.CallWithContext(ctx, m, 0, zone, portProto).Err
}

func (c *FirewallDebusClient) AddService(ctx context.Context, zone string, service string) error {
	m := prefix + "zone.addService"
	return c.obj.CallWithContext(ctx, m, 0, zone, service, 0).Err
}

func (c *FirewallDebusClient) RemoveService(ctx context.Context, zone string, service string) error {
	m := prefix + "zone.removeService"
	return c.obj.CallWithContext(ctx, m, 0, zone, service).Err
}
