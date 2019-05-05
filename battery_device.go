package hubitat

import (
	"encoding/json"
	"fmt"
)

// BatteryDevice represents a device that can act as a battery.
type BatteryDevice struct {
	Device
}

// Battery returns the battery percentage of the device.
func (d *BatteryDevice) Battery() (float64, error) {
	if battery, ok := d.Attributes["battery"].(json.Number); ok {
		v, err := battery.Float64()

		if err != nil {
			return 0, fmt.Errorf("incorrect battery attribute in device %s", d)
		}

		return v / 100, nil
	}

	return 0, fmt.Errorf("missing `battery` attribute in device %s", d)
}
