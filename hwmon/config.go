package hwmon

import (
	"encoding/json"
	"errors"
	"log"

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

type Sensor struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Name        *[]string `json:"driver_hwmon"`
	NameNew     *string   `json:"name"`
}

type ISensor interface {
	Collect(metrics chan<- prometheus.Metric)
}

type VoltageSensor struct {
  Min         *float64 `json:"min"`
  Max         *float64 `json:"max"`
  CriticalMin *float64 `json:"lcrit"`
  CriticalMax *float64 `json:"crit"`
  Input       *float64 `json:"input"`
	Sensor
}

func (s *VoltageSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, voltageMinDesc,         s.Min)
	s.CollectNum(metrics, voltageMaxDesc,         s.Max)
	s.CollectNum(metrics, voltageCriticalMinDesc, s.CriticalMin)
	s.CollectNum(metrics, voltageCriticalMaxDesc, s.CriticalMax)
	s.CollectNum(metrics, voltageDesc,            s.Input)
}

type FanSensor struct {
  Min    *float64 `json:"min"`
  Max    *float64 `json:"max"`
  Input  *float64 `json:"input"`
  Pulses *float64 `json:"pulses"`
  Target *float64 `json:"target"`
	Sensor
}

func (s *FanSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, fanMinDesc,    s.Min)
	s.CollectNum(metrics, fanMaxDesc,    s.Max)
	s.CollectNum(metrics, fanDesc,       s.Input)
	s.CollectNum(metrics, fanPulsesDesc, s.Pulses)
	s.CollectNum(metrics, fanTargetDesc, s.Target)
}

type TemperatureSensor struct {
  Min              *float64 `json:"min"`
  Max              *float64 `json:"max"`
  MinHysteresis    *float64 `json:"min_hyst"`
  MaxHysteresis    *float64 `json:"max_hyst"`
  Input            *float64 `json:"input"`
  CriticalMin      *float64 `json:"lcrit"`
  CriticalMax      *float64 `json:"crit"`
  CriticalMinHyst  *float64 `json:"lcrit_hyst"`
  CriticalMaxHyst  *float64 `json:"crit_hyst"`
  EmergencyMax     *float64 `json:"emergency"`
  EmergencyMaxHyst *float64 `json:"emergency_hyst"`
  Offset           *float64 `json:"offset"`
	Sensor
}

func (s *TemperatureSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, temperatureMinDesc,                    s.Min)
	s.CollectNum(metrics, temperatureMaxDesc,                    s.Max)
	s.CollectNum(metrics, temperatureMinHysteresisDesc,          s.MinHysteresis)
	s.CollectNum(metrics, temperatureMaxHysteresisDesc,          s.MaxHysteresis)
	s.CollectNum(metrics, temperatureDesc,                       s.Input)
	s.CollectNum(metrics, temperatureCriticalMinDesc,            s.CriticalMin)
	s.CollectNum(metrics, temperatureCriticalMaxDesc,            s.CriticalMax)
	s.CollectNum(metrics, temperatureCriticalMinHysteresisDesc,  s.CriticalMinHyst)
	s.CollectNum(metrics, temperatureCriticalMaxHysteresisDesc,  s.CriticalMaxHyst)
	s.CollectNum(metrics, temperatureEmergencyMaxDesc,           s.EmergencyMax)
	s.CollectNum(metrics, temperatureEmergencyMaxHysteresisDesc, s.EmergencyMaxHyst)
	s.CollectNum(metrics, temperatureOffsetDesc,                 s.Offset)
}

type CurrentSensor struct {
  Min         *float64 `json:"min"`
  Max         *float64 `json:"max"`
  CriticalMin *float64 `json:"lcrit"`
  CriticalMax *float64 `json:"crit"`
  Input       *float64 `json:"input"`
	Sensor
}

func (s *CurrentSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, currentMinDesc,          s.Min)
  s.CollectNum(metrics, currentMaxDesc,          s.Max)
  s.CollectNum(metrics, currentCriticalMinValue, s.CriticalMin)
  s.CollectNum(metrics, currentCriticalMaxValue, s.CriticalMax)
  s.CollectNum(metrics, currentDesc,             s.Input)
}

type PowerSensor struct {
	State string `json:"state"`
	Sensor
}

// func (s *PowerSensor) UnmarshalJSON(b []byte) error {
// 	type powerSensor PowerSensor

// 	err := json.Unmarshal(b, (*powerSensor)(s))
// 	if err != nil {
// 		return err
// 	}

// 	var data map[string]interface{};
// 	err = json.Unmarshal(b, &data)
// 	if err != nil {
// 		return err
// 	}

// 	for key, value := range data {
// 		if strings.Contains(key, "present") || strings.Contains(key, "status") {

// 		}

// 		if strings.Contains(key, "all_ok") {

// 		}
// 	}

// 	strings.Contains()
// }

func (s *PowerSensor) Collect(metrics chan<- prometheus.Metric) {
	isOk := s.State == "OK"
	s.CollectBool(metrics, powerAllOk, &isOk)
}

type RawSensor struct {
  Raw *float64 `json:"raw"`
	Sensor
}

func (s *RawSensor) Collect(metrics chan<- prometheus.Metric) {
	s.CollectNum(metrics, rawValueDesc, s.Raw)
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

		log.Println(sensor)

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
