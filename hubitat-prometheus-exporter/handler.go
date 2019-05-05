package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/briandowns/openweathermap"
	"github.com/intelux/hubitat"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var prometheusHandler = promhttp.Handler()

const (
	deviceLabel = `device`
)

var (
	hubitatBatteryLevelCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_battery_level_current",
			Help: "The current battery level percentage of devices.",
		},
		[]string{deviceLabel},
	)
	hubitatTemperatureCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_temperature_current",
			Help: "The current temperature of devices.",
		},
		[]string{deviceLabel},
	)
	hubitatIlluminanceCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_illuminance_current",
			Help: "The current illuminance of devices.",
		},
		[]string{deviceLabel},
	)
	hubitatHumidityCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_humidity_current",
			Help: "The current humidity percentage of devices.",
		},
		[]string{deviceLabel},
	)
	hubitatSwitchCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_switch_current",
			Help: "The current switch state of devices.",
		},
		[]string{deviceLabel},
	)
	hubitatSwitchLevelCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hubitat_switch_level_current",
			Help: "The current switch level of devices.",
		},
		[]string{deviceLabel},
	)
)

// Handler implements a handler.
type Handler struct {
	Client             *hubitat.Client
	CurrentWeatherData *openweathermap.CurrentWeatherData
	CityID             int
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	devices, err := h.Client.GetDevices(req.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s\n", err)

		return
	}

	hubitatBatteryLevelCurrent.Reset()
	hubitatTemperatureCurrent.Reset()
	hubitatIlluminanceCurrent.Reset()
	hubitatHumidityCurrent.Reset()
	hubitatSwitchCurrent.Reset()
	hubitatSwitchLevelCurrent.Reset()

	for _, device := range devices.BatteryDevices() {
		value, err := device.Battery()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get battery level for `%s`: %s", device, err)
		} else {
			hubitatBatteryLevelCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	for _, device := range devices.TemperatureDevices() {
		value, err := device.Temperature()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get temperature for `%s`: %s", device, err)
		} else {
			hubitatTemperatureCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	for _, device := range devices.IlluminanceDevices() {
		value, err := device.Illuminance()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get illuminance for `%s`: %s", device, err)
		} else {
			hubitatIlluminanceCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	for _, device := range devices.HumidityDevices() {
		value, err := device.Humidity()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get humidity for `%s`: %s", device, err)
		} else {
			hubitatHumidityCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	for _, device := range devices.SwitchDevices() {
		onOff, err := device.Switch()

		value := 0.0

		if onOff {
			value = 1.0
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get switch for `%s`: %s", device, err)
		} else {
			hubitatSwitchCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	for _, device := range devices.SwitchLevelDevices() {
		value, err := device.SwitchLevel()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get switch level for `%s`: %s", device, err)
		} else {
			hubitatSwitchLevelCurrent.With(prometheus.Labels{
				deviceLabel: device.String(),
			}).Set(value)
		}
	}

	if h.CurrentWeatherData != nil {
		if err = h.CurrentWeatherData.CurrentByID(h.CityID); err == nil {
			hubitatTemperatureCurrent.With(prometheus.Labels{
				deviceLabel: "Outside",
			}).Set(h.CurrentWeatherData.Main.Temp)
			hubitatHumidityCurrent.With(prometheus.Labels{
				deviceLabel: "Outside",
			}).Set(float64(h.CurrentWeatherData.Main.Humidity) / 100.0)
		}
	}

	prometheusHandler.ServeHTTP(w, req)
}
