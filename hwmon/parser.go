package hwmon

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func (s *Sensor) GetName() *string {
	if s.NameNew != nil {
		return s.NameNew
	}

	if s.Name != nil && len(*s.Name) != 0 {
		return &(*s.Name)[0]
	}

	return nil
}

func (s *Sensor) CollectNum(metrics chan<- prometheus.Metric, desc *prometheus.Desc, value *float64) {
	if value == nil {
		return
	}

	name := s.GetName()
	if name == nil {
		return
	}

	metric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, *value, s.Description, *name)
	metrics <- metric
}

func (s *Sensor) CollectBool(metrics chan<- prometheus.Metric, desc *prometheus.Desc, value *bool) {
	if value == nil {
		return
	}

	name := s.GetName()
	if name == nil {
		return
	}

	num := 0.0
	if *value {
		num = 1.0
	}

	metric := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, num, s.Description, *name)
	metrics <- metric
}

// Because of funny golang unmarshalling and prometheus metric safety this is necessary
// to unmarshal to either a float64 or to a float64 from a string or nil
type OptFloatString struct {
	inner *float64
}

func (opt *OptFloatString) UnmarshalJSON(b []byte) error {
	var inner any
	
	err := json.Unmarshal(b, (&inner))
	if err != nil {
		return err
	}

	switch v := inner.(type) {
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		opt.inner = &parsed
	case float64:
		opt.inner = &v
	case int:
		intermediate := float64(v)
		opt.inner = &intermediate
	}

	return nil
}



type Sensor struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Name        *[]string `json:"driver_hwmon,omitempty"`
	NameNew     *string   `json:"name,omitempty"`
}

type ISensor interface {
	Collect(metrics chan<- prometheus.Metric)
}



type VoltageSensor struct {
  Min         OptFloatString `json:"min,omitempty"`
  Max         OptFloatString `json:"max,omitempty"`
  CriticalMin OptFloatString `json:"lcrit,omitempty"`
  CriticalMax OptFloatString `json:"crit,omitempty"`
  Input       OptFloatString `json:"input,omitempty"`
	Sensor
}

func (s *VoltageSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, voltageMinDesc,         s.Min.inner)
	s.CollectNum(metrics, voltageMaxDesc,         s.Max.inner)
	s.CollectNum(metrics, voltageCriticalMinDesc, s.CriticalMin.inner)
	s.CollectNum(metrics, voltageCriticalMaxDesc, s.CriticalMax.inner)
	s.CollectNum(metrics, voltageDesc,            s.Input.inner)
}



type FanSensor struct {
  Min    OptFloatString `json:"min,omitempty"`
  Max    OptFloatString `json:"max,omitempty"`
  Input  OptFloatString `json:"input,omitempty"`
  Pulses OptFloatString `json:"pulses,omitempty"`
  Target OptFloatString `json:"target,omitempty"`
	Sensor
}

func (s *FanSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, fanMinDesc,    s.Min.inner)
	s.CollectNum(metrics, fanMaxDesc,    s.Max.inner)
	s.CollectNum(metrics, fanDesc,       s.Input.inner)
	s.CollectNum(metrics, fanPulsesDesc, s.Pulses.inner)
	s.CollectNum(metrics, fanTargetDesc, s.Target.inner)
}



type TemperatureSensor struct {
  Min              OptFloatString `json:"min,omitempty"`
  Max              OptFloatString `json:"max,omitempty"`
  MinHysteresis    OptFloatString `json:"min_hyst,omitempty"`
  MaxHysteresis    OptFloatString `json:"max_hyst,omitempty"`
  Input            OptFloatString `json:"input,omitempty"`
  CriticalMin      OptFloatString `json:"lcrit,omitempty"`
  CriticalMax      OptFloatString `json:"crit,omitempty"`
  CriticalMinHyst  OptFloatString `json:"lcrit_hyst,omitempty"`
  CriticalMaxHyst  OptFloatString `json:"crit_hyst,omitempty"`
  EmergencyMax     OptFloatString `json:"emergency,omitempty"`
  EmergencyMaxHyst OptFloatString `json:"emergency_hyst,omitempty"`
  Offset           OptFloatString `json:"offset,omitempty"`
	Sensor
}

func (s *TemperatureSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, temperatureMinDesc,                    s.Min.inner)
	s.CollectNum(metrics, temperatureMaxDesc,                    s.Max.inner)
	s.CollectNum(metrics, temperatureMinHysteresisDesc,          s.MinHysteresis.inner)
	s.CollectNum(metrics, temperatureMaxHysteresisDesc,          s.MaxHysteresis.inner)
	s.CollectNum(metrics, temperatureDesc,                       s.Input.inner)
	s.CollectNum(metrics, temperatureCriticalMinDesc,            s.CriticalMin.inner)
	s.CollectNum(metrics, temperatureCriticalMaxDesc,            s.CriticalMax.inner)
	s.CollectNum(metrics, temperatureCriticalMinHysteresisDesc,  s.CriticalMinHyst.inner)
	s.CollectNum(metrics, temperatureCriticalMaxHysteresisDesc,  s.CriticalMaxHyst.inner)
	s.CollectNum(metrics, temperatureEmergencyMaxDesc,           s.EmergencyMax.inner)
	s.CollectNum(metrics, temperatureEmergencyMaxHysteresisDesc, s.EmergencyMaxHyst.inner)
	s.CollectNum(metrics, temperatureOffsetDesc,                 s.Offset.inner)
}



type CurrentSensor struct {
  Min         OptFloatString `json:"min,omitempty"`
  Max         OptFloatString `json:"max,omitempty"`
  CriticalMin OptFloatString `json:"lcrit,omitempty"`
  CriticalMax OptFloatString `json:"crit,omitempty"`
  Input       OptFloatString `json:"input,omitempty"`
	Sensor
}

func (s *CurrentSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, currentMinDesc,          s.Min.inner)
  s.CollectNum(metrics, currentMaxDesc,          s.Max.inner)
  s.CollectNum(metrics, currentCriticalMinValue, s.CriticalMin.inner)
  s.CollectNum(metrics, currentCriticalMaxValue, s.CriticalMax.inner)
  s.CollectNum(metrics, currentDesc,             s.Input.inner)
}



// {
//   "name": "PSU1",
//   "start_time": 1725878482,
//   "psu1_pwr_status": 1, # Is power present
//   "state": "OK",        # We also check this for is all ok
//   "prev_state": "OK",   # And this
//   "prev_msg": null,
//   "msg": null,
//   "psu1_power": 35,     # Power value
//   "log_time": 1725878482,
//   "type": "power",
//   "psu1_status": 1,
//   "description": "PSU1"
// },
//
// or
//
// # No power value
// {
//   "driver_hwmon": [
//       "psu_pwr2"
//   ],
//   "name": "PSU2",
//   "start_time": 1745417928,
//   "psu_pwr2_all_ok": "1",  # Is all ok
//   "state": "OK",           # But we check this instead
//   "prev_state": "OK",      # Last state
//   "psu_pwr2_present": "1", # Is power present
//   "driver_path": "/sys/bus/i2c/devices/0-0030",
//   "msg": null,
//   "prev_msg": null,
//   "log_time": 1745417928,
//   "type": "power",
//   "description": "PSU2"
// },

type PowerSensor struct {
	State     string   `json:"state"`
	PrevState string   `json:"prev_state"`
	IsPresent *bool    `json:"-"`
	Power     *float64 `json:"-"` // Not always available
	Sensor
}

func (s *PowerSensor) UnmarshalJSON(b []byte) error {
	type powerSensor PowerSensor

	err := json.Unmarshal(b, (*powerSensor)(s))
	if err != nil {
		return err
	}

	var data map[string]json.RawMessage;
	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	for key, value := range data {
		var opt OptFloatString
		err := json.Unmarshal(value, &opt)
		if err != nil || opt.inner == nil {
			continue
		}
		num := *opt.inner

		if strings.Contains(key, "present") || strings.Contains(key, "pwr_status") {
			isPresent := num > 0
			s.IsPresent = &isPresent
			continue
		}

		if strings.Contains(key, "power") {
			s.Power = &num
		}
	}

	return nil
}

func OkValueToFloat(state string) (float64, error) {
	switch state {
	case "OK":
		return 1, nil
	case "BAD":
		return 0, nil
	case "POWERED OFF":
		return -1, nil
	case "NOT DETECTED":
		return -2, nil
	}
	return float64(-1), fmt.Errorf("could not parse state value of %s to a float", state)
}

func (s *PowerSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, powerWatt, s.Power)
	s.CollectBool(metrics, powerPresent, s.IsPresent)

	stateOkF, err := OkValueToFloat(s.State)
	if err == nil {
		s.CollectNum(metrics, powerAllOk, &stateOkF)
	}

	stateOkPrevF, err := OkValueToFloat(s.PrevState)
	if err == nil {
		s.CollectNum(metrics, powerAllOkPrev, &stateOkPrevF)
	}
}

type RawSensor struct {
  Raw OptFloatString `json:"raw"`
	Sensor
}

func (s *RawSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, rawValueDesc, s.Raw.inner)
}

func UnmarshalSensors(data []byte) ([]ISensor, error) {
	var rawSensors []json.RawMessage

	err := json.Unmarshal(data, &rawSensors)
	if err != nil {
		return nil, err
	}

	var sensors []ISensor
	
	for _, raw := range rawSensors {
		var sensor Sensor

		err = json.Unmarshal(raw, &sensor)
		if err != nil {
			return nil, err
		}

		var i ISensor

		switch sensor.Type {
		case "voltage":
			i = &VoltageSensor{}
		case "fan":
			i = &FanSensor{}
		case "temp":
			i = &TemperatureSensor{}
		case "current":
			i = &CurrentSensor{}
		case "power":
			i = &PowerSensor{}
		case "raw":
			i = &RawSensor{}
		default:
			return nil, errors.New("unknown sensor type")
		}

		err = json.Unmarshal(raw, i)
		if err != nil {
			return nil, err
		}

		sensors = append(sensors, i)
	}

	return sensors, nil
}
