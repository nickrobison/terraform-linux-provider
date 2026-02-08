package firewall

import (
	"github.com/godbus/dbus/v5"
	"github.com/nickrobison/terraform-linux-provider/server/bus"
	"github.com/rs/zerolog"
)

type ZoneObject struct {
	obj    dbus.BusObject
	logger *zerolog.Logger
}

func (o *ZoneObject) Name() (string, error) {
	property := prefix + "zone.name"
	name, err := bus.Decode[string](o.logger, o.obj, property)
	return name, err
}

func (o *ZoneObject) Description() (string, error) {
	property := prefix + "zone.description"
	description, err := bus.Decode[string](o.logger, o.obj, property)
	return description, err
}

func (o *ZoneObject) Target() (string, error) {
	property := prefix + "zone.target"
	target, err := bus.Decode[string](o.logger, o.obj, property)
	return target, err
}

func (o *ZoneObject) Services() ([]string, error) {
	property := prefix + "zone.services"
	services, err := bus.Decode[[]string](o.logger, o.obj, property)
	return services, err
}

func (o *ZoneObject) Ports() ([]string, error) {
	property := prefix + "zone.ports"
	ports, err := bus.Decode[[]string](o.logger, o.obj, property)
	return ports, err
}

func (o *ZoneObject) RichRules() ([]string, error) {
	property := prefix + "zone.rich_rules"
	rules, err := bus.Decode[[]string](o.logger, o.obj, property)
	return rules, err
}

func NewZoneObject(obj dbus.BusObject, logger *zerolog.Logger) *ZoneObject {
	log := logger.With().Str("path", string(obj.Path())).Logger()
	return &ZoneObject{
		obj:    obj,
		logger: &log,
	}
}
