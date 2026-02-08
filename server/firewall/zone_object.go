package firewall

import (
"github.com/godbus/dbus/v5"
"github.com/rs/zerolog"
)

type ZoneObject struct {
obj      dbus.BusObject
logger   *zerolog.Logger
name     string
settings map[string]interface{}
}

func (o *ZoneObject) Name() (string, error) {
return o.name, nil
}

func (o *ZoneObject) Description() (string, error) {
if o.settings == nil {
if err := o.loadSettings(); err != nil {
return "", err
}
}
if desc, ok := o.settings["description"].(string); ok {
return desc, nil
}
return "", nil
}

func (o *ZoneObject) Target() (string, error) {
if o.settings == nil {
if err := o.loadSettings(); err != nil {
return "", err
}
}
if target, ok := o.settings["target"].(string); ok {
return target, nil
}
return "", nil
}

func (o *ZoneObject) Services() ([]string, error) {
if o.settings == nil {
if err := o.loadSettings(); err != nil {
return nil, err
}
}
if services, ok := o.settings["services"].([]string); ok {
return services, nil
}
return []string{}, nil
}

func (o *ZoneObject) Ports() ([]string, error) {
if o.settings == nil {
if err := o.loadSettings(); err != nil {
return nil, err
}
}
// Ports are returned as array of [port, protocol] tuples
// Convert to "port/protocol" string format
var result []string
if ports, ok := o.settings["ports"].([]interface{}); ok {
for _, p := range ports {
if portTuple, ok := p.([]interface{}); ok && len(portTuple) >= 2 {
port := portTuple[0].(string)
protocol := portTuple[1].(string)
result = append(result, port+"/"+protocol)
}
}
}
return result, nil
}

func (o *ZoneObject) RichRules() ([]string, error) {
if o.settings == nil {
if err := o.loadSettings(); err != nil {
return nil, err
}
}
if rules, ok := o.settings["rich_rules"].([]string); ok {
return rules, nil
}
return []string{}, nil
}

func (o *ZoneObject) loadSettings() error {
// Use getSettings method to retrieve zone configuration
// Returns a tuple with all zone settings
method := "org.fedoraproject.FirewallD1.config.zone.getSettings"
call := o.obj.Call(method, 0)
if call.Err != nil {
return call.Err
}

// Parse the settings tuple returned by getSettings
// Format: (version, name, description, unused, target, [services], [ports], ...)
var settingsTuple []interface{}
if err := call.Store(&settingsTuple); err != nil {
return err
}

// Extract relevant fields from the tuple
o.settings = make(map[string]interface{})
if len(settingsTuple) > 2 {
o.settings["description"] = settingsTuple[2]
}
if len(settingsTuple) > 4 {
o.settings["target"] = settingsTuple[4]
}
if len(settingsTuple) > 5 {
o.settings["services"] = settingsTuple[5]
}
if len(settingsTuple) > 6 {
o.settings["ports"] = settingsTuple[6]
}
if len(settingsTuple) > 12 {
o.settings["rich_rules"] = settingsTuple[12]
}

return nil
}

func NewZoneObject(obj dbus.BusObject, logger *zerolog.Logger, name string) *ZoneObject {
log := logger.With().Str("path", string(obj.Path())).Logger()
return &ZoneObject{
obj:    obj,
logger: &log,
name:   name,
}
}
