package hubitat

const (
	// CapabilityBattery represents all the devices that can act as a battery.
	CapabilityBattery Capability = "Battery"
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
