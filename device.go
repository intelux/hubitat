package hubitat

import (
	"time"
)

// DeviceID represents a device identifier.
type DeviceID int

// Device represents a device, as returned by the Maker API.
type Device struct {
	Name         string                 `json:"name"`
	Label        string                 `json:"label,omitempty"`
	Type         string                 `json:"type"`
	ID           DeviceID               `json:"id"`
	Date         time.Time              `json:"date"`
	Model        *string                `json:"model"`
	Manufacturer *string                `json:"manufacturer"`
	Capabilities Capabilities           `json:"capabilities"`
	Attributes   map[string]interface{} `json:"attributes"`
	Commands     []Command              `json:"commands"`
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
