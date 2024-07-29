package bus

import (
	"github.com/godbus/dbus/v5"
	"github.com/rs/zerolog"
)

func Decode[T any](log *zerolog.Logger, obj dbus.BusObject, property string) (T, error) {
	var v T
	prop, err := obj.GetProperty(property)
	if err != nil {
		return v, err
	}
	log.Debug().Str("property", property).Str("variant", prop.Signature().String()).Msg("Received property")
	err = prop.Store(&v)
	if err != nil {
		return v, err
	}
	return v, nil
}
