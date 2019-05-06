package hubitat

const (
	// CapabilityBattery represents all the devices that can act as a battery.
	CapabilityBattery Capability = "Battery"
	// CapabilityTemperatureMeasurement represents all the devices can measure temperature.
	CapabilityTemperatureMeasurement = "TemperatureMeasurement"
	// CapabilityIlluminanceMeasurement represents all the devices can measure illuminance.
	CapabilityIlluminanceMeasurement = "IlluminanceMeasurement"
	// CapabilityRelativeHumidityMeasurement represents all the devices can measure humidity.
	CapabilityRelativeHumidityMeasurement = "RelativeHumidityMeasurement"
	// CapabilitySwitch represents all the devices can act as a switch.
	CapabilitySwitch = "Switch"
	// CapabilitySwitchLevel represents all the devices can act as a switch with level.
	CapabilitySwitchLevel = "SwitchLevel"
	// CapabilityLock represents all the devices can act as locks.
	CapabilityLock = "Lock"
)

// Contains checks whether a capability is contained in the list of capabilities.
func (c Capabilities) Contains(capability Capability) bool {
	for _, capa := range c {
		if capa == capability {
			return true
		}
	}

	return false
}

// Capability represents a device capability.
type Capability string

// Capabilities represents a list of capabilities.
type Capabilities []Capability
