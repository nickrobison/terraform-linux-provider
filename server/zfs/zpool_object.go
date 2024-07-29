package zfs

import (
	"github.com/godbus/dbus/v5"
	"github.com/nickrobison/terraform-linux-provider/server/bus"
	"github.com/rs/zerolog"
)

type ZpoolObject struct {
	obj    dbus.BusObject
	logger *zerolog.Logger
}

func (o ZpoolObject) Name() (string, error) {
	property := prefix + "Pool.Name"
	name, err := bus.Decode[string](o.logger, o.obj, property)
	return name, err
}

func NewZpoolObject(obj dbus.BusObject, logger *zerolog.Logger) *ZpoolObject {
	log := logger.With().Str("path", string(obj.Path())).Logger()
	return &ZpoolObject{
		obj:    obj,
		logger: &log,
	}

}
