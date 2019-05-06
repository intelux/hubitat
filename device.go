package hubitat

import (
	"fmt"
	"strconv"
)

// DeviceID represents a device identifier.
type DeviceID string

// Device represents a device, as returned by the Maker API.
type Device struct {
	Name         string       `json:"name"`
	Label        string       `json:"label,omitempty"`
	Type         string       `json:"type"`
	ID           DeviceID     `json:"id"`
	Model        *string      `json:"model"`
	Manufacturer *string      `json:"manufacturer"`
	Capabilities Capabilities `json:"capabilities"`
	Attributes   Attributes   `json:"attributes"`
	Commands     []Command    `json:"commands"`
}

func (d Device) String() string {
	if d.Label != "" {
		return d.Label
	}

	return d.Name
}

// Devices represents a list of devices.
type Devices []Device

// BatteryDevices returns all the battery devices.
func (d Devices) BatteryDevices() (devices []BatteryDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilityBattery) {
			devices = append(devices, BatteryDevice{device})
		}
	}

	return
}

// TemperatureDevices returns all the temperature devices.
func (d Devices) TemperatureDevices() (devices []TemperatureDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilityTemperatureMeasurement) {
			devices = append(devices, TemperatureDevice{device})
		}
	}

	return
}

// IlluminanceDevices returns all the temperature devices.
func (d Devices) IlluminanceDevices() (devices []IlluminanceDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilityIlluminanceMeasurement) {
			devices = append(devices, IlluminanceDevice{device})
		}
	}

	return
}

// HumidityDevices returns all the temperature devices.
func (d Devices) HumidityDevices() (devices []HumidityDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilityRelativeHumidityMeasurement) {
			devices = append(devices, HumidityDevice{device})
		}
	}

	return
}

// SwitchDevices returns all the switch devices.
func (d Devices) SwitchDevices() (devices []SwitchDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilitySwitch) {
			devices = append(devices, SwitchDevice{device})
		}
	}

	return
}

// SwitchLevelDevices returns all the switch level devices.
func (d Devices) SwitchLevelDevices() (devices []SwitchLevelDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilitySwitchLevel) {
			devices = append(devices, SwitchLevelDevice{device})
		}
	}

	return
}

// LockDevices returns all the lock devices.
func (d Devices) LockDevices() (devices []LockDevice) {
	for _, device := range d {
		if device.Capabilities.Contains(CapabilityLock) {
			devices = append(devices, LockDevice{device})
		}
	}

	return
}

// Attributes represents device attributes.
type Attributes map[string]interface{}

// GetNumber returns a given value as a number.
func (a Attributes) GetNumber(key string) (float64, error) {
	v, ok := a[key]

	if !ok {
		return 0, fmt.Errorf("attributes contains no such key: %s", key)
	}

	s, ok := v.(string)

	if !ok {
		return 0, fmt.Errorf("attributes value %s was exptected to be a string", key)
	}

	f, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0, fmt.Errorf("attribute value for %s is not a valid float: %s", key, err)
	}

	return f, nil
}

// GetBoolean returns a given value as a boolean.
func (a Attributes) GetBoolean(key string, trueValue string, falseValue string) (bool, error) {
	v, ok := a[key]

	if !ok {
		return false, fmt.Errorf("attributes contains no such key: %s", key)
	}

	s, ok := v.(string)

	if !ok {
		return false, fmt.Errorf("attributes value %s was exptected to be a string", key)
	}

	switch s {
	case trueValue:
		return true, nil
	case falseValue:
		return false, nil
	default:
		return false, fmt.Errorf("attribute value for %s is not valid (expected `%s` or `%s`, got `%s`)", key, trueValue, falseValue, s)
	}
}

// GetPercentage returns a given value as a percentage.
func (a Attributes) GetPercentage(key string) (f float64, err error) {
	f, err = a.GetNumber(key)
	f /= 100

	return
}

// BatteryDevice represents a device that can act as a battery.
type BatteryDevice struct {
	Device
}

// Battery returns the battery percentage of the device.
func (d *BatteryDevice) Battery() (float64, error) {
	return d.Attributes.GetPercentage("battery")
}

// TemperatureDevice represents a device that can read temperature.
type TemperatureDevice struct {
	Device
}

// Temperature returns the temperature, in the current device unit.
func (d *TemperatureDevice) Temperature() (float64, error) {
	return d.Attributes.GetNumber("temperature")
}

// IlluminanceDevice represents a device that can read temperature.
type IlluminanceDevice struct {
	Device
}

// Illuminance returns the temperature, in the current device unit.
func (d *IlluminanceDevice) Illuminance() (float64, error) {
	return d.Attributes.GetNumber("illuminance")
}

// HumidityDevice represents a device that can read humidity.
type HumidityDevice struct {
	Device
}

// Humidity returns the humidity, in percent.
func (d *HumidityDevice) Humidity() (float64, error) {
	return d.Attributes.GetPercentage("temperature")
}

// SwitchDevice represents a device that can read humidity.
type SwitchDevice struct {
	Device
}

// Switch returns the switch level, in percent.
func (d *SwitchDevice) Switch() (bool, error) {
	return d.Attributes.GetBoolean("switch", "on", "off")
}

// SwitchLevelDevice represents a device that can read humidity.
type SwitchLevelDevice struct {
	Device
}

// SwitchLevel returns the switch level, in percent.
func (d *SwitchLevelDevice) SwitchLevel() (float64, error) {
	return d.Attributes.GetPercentage("level")
}

// LockDevice represents a device that can be locked.
type LockDevice struct {
	Device
}

// Lock returns the locked state of the device.
func (d *LockDevice) Lock() (bool, error) {
	return d.Attributes.GetBoolean("lock", "locked", "unlocked")
}
