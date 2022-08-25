package hwmon

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wobcom/cumulus-exporter/util"
	"sync"
)

const prefix = "hwmon_"

var (
	voltageMinDesc           *prometheus.Desc
	voltageCriticalMinDesc   *prometheus.Desc
	voltageMaxDesc           *prometheus.Desc
	voltageCriticalMaxDesc   *prometheus.Desc
	voltageDesc              *prometheus.Desc
	voltageLabelInfoDesc     *prometheus.Desc
	voltageSensorEnabledDesc *prometheus.Desc

	fanMinDesc           *prometheus.Desc
	fanMaxDesc           *prometheus.Desc
	fanDesc              *prometheus.Desc
	fanDivisorDesc       *prometheus.Desc
	fanPulsesDesc        *prometheus.Desc
	fanTargetDesc        *prometheus.Desc
	fanLabelsDesc        *prometheus.Desc
	fanSensorEnabledDesc *prometheus.Desc

	temperatureTypeDesc                   *prometheus.Desc
	temperatureMaxDesc                    *prometheus.Desc
	temperatureMinDesc                    *prometheus.Desc
	temperatureMaxHysteresisDesc          *prometheus.Desc
	temperatureMinHysteresisDesc          *prometheus.Desc
	temperatureDesc                       *prometheus.Desc
	temperatureCriticalMaxDesc            *prometheus.Desc
	temperatureCriticalMaxHysteresisDesc  *prometheus.Desc
	temperatureEmergencyMaxDesc           *prometheus.Desc
	temperatureEmergencyMaxHysteresisDesc *prometheus.Desc
	temperatureCriticalMinDesc            *prometheus.Desc
	temperatureCriticalMinHysteresisDesc  *prometheus.Desc
	temperatureOffsetDesc                 *prometheus.Desc
	temperatureLabelDesc                  *prometheus.Desc
	temperatureSensorEnabledDesc          *prometheus.Desc

	currentMaxDesc           *prometheus.Desc
	currentMinDesc           *prometheus.Desc
	currentCriticalMinValue  *prometheus.Desc
	currentCriticalMaxValue  *prometheus.Desc
	currentDesc              *prometheus.Desc
	currentSensorEnabledDesc *prometheus.Desc

	rawValueDesc *prometheus.Desc
)

func init() {
	sensorLabels := []string{"driver_path", "hw_mon", "description"}
	channelLabels := []string{"driver_path", "hw_mon", "description", "channel"}

	voltageMinDesc = prometheus.NewDesc(prefix+"voltage_min_volts", "Voltage min value. Unit: Volts", sensorLabels, nil)
	voltageCriticalMinDesc = prometheus.NewDesc(prefix+"voltage_critical_min_volts", "Voltage critical min value. Unit: Volts", sensorLabels, nil)
	voltageMaxDesc = prometheus.NewDesc(prefix+"voltage_max_volts", "Voltage max value. Unit: Volts", sensorLabels, nil)
	voltageCriticalMaxDesc = prometheus.NewDesc(prefix+"voltage_critical_max_volts", "Voltage critical max value. Unit: Volts", sensorLabels, nil)
	voltageDesc = prometheus.NewDesc(prefix+"voltage_volts", "Voltage input value. Unit: Volts", sensorLabels, nil)
	voltageLabelInfoDesc = prometheus.NewDesc(prefix+"voltage_info", "Suggested voltage channel label.", channelLabels, nil)
	voltageSensorEnabledDesc = prometheus.NewDesc(prefix+"voltage_sensor_enabled_bool", "1 = sensor enabled, 0 = sensor disabled", sensorLabels, nil)

	fanMinDesc = prometheus.NewDesc(prefix+"fan_min_rpm", "Fan minimum value. Unit: revolution/min", sensorLabels, nil)
	fanMaxDesc = prometheus.NewDesc(prefix+"fan_max_rpm", "Fan maximum value. Unit: revolution/min", sensorLabels, nil)
	fanDesc = prometheus.NewDesc(prefix+"fan_rpm", "Fan input value. Unit: revolution/min", sensorLabels, nil)
	fanDivisorDesc = prometheus.NewDesc(prefix+"fan_divisor", "Fan divisor. Integer value in powers of 2 (1, 2, 4, 8, 16, 32, 64, 128).", sensorLabels, nil)
	fanPulsesDesc = prometheus.NewDesc(prefix+"fan_pulses", "Number of tachometer pulses per fan revolution", sensorLabels, nil)
	fanTargetDesc = prometheus.NewDesc(prefix+"fan_target_rpm", "Desired fan speed. Unit: revolution/min", sensorLabels, nil)
	fanLabelsDesc = prometheus.NewDesc(prefix+"fan_info", "Suggested fan channel label", channelLabels, nil)
	fanSensorEnabledDesc = prometheus.NewDesc(prefix+"fan_sensor_enabled_bool", "1 = sensor enabled, 0 = sensor disabled", sensorLabels, nil)

	temperatureTypeDesc = prometheus.NewDesc(prefix+"temperature_sensor_type_selection_info", "Sensor type selection.", append(sensorLabels, "sensor_type"), nil)
	temperatureMaxDesc = prometheus.NewDesc(prefix+"temperature_max_celsius", "Temperature max value. Unit: degree Celsius", sensorLabels, nil)
	temperatureMinDesc = prometheus.NewDesc(prefix+"temperature_min_celsius", "Temperature min value. Unit: degree Celsius", sensorLabels, nil)
	temperatureMaxHysteresisDesc = prometheus.NewDesc(prefix+"temperature_max_hysteresis_celsius", "Temperature hysteresis value for max limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureMinHysteresisDesc = prometheus.NewDesc(prefix+"temperature_min_hysteresis_celsius", "Temperature hysteresis value for min limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureDesc = prometheus.NewDesc(prefix+"temperature_celsius", "Temperature input value. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMaxDesc = prometheus.NewDesc(prefix+"temperature_critical_max_celsius", "Temperature critical max value, typically greater than corresponding temp_max values. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMaxHysteresisDesc = prometheus.NewDesc(prefix+"temperature_critical_max_hysteresis_celsius", "Temperature hysteresis value for critical limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureEmergencyMaxDesc = prometheus.NewDesc(prefix+"temperature_emergency_max_celsius", "Temperature emergency max value, for chips supporting more than two upper temperature limits. Unit: degree Celsius", sensorLabels, nil)
	temperatureEmergencyMaxHysteresisDesc = prometheus.NewDesc(prefix+"temperature_emergency_max_hysteresis_celsius", "Temperature hysteresis value for emergency limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMinDesc = prometheus.NewDesc(prefix+"temperature_critical_min_celsius", "Temperature criticial min value, typically lower than corresponding temp_min values. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMinHysteresisDesc = prometheus.NewDesc(prefix+"temperature_critical_min_hysteresis_celsius", "Temperature hysteresis value for critical min limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureOffsetDesc = prometheus.NewDesc(prefix+"temperature_offset_celsius", "Temperature offset which is added to the temperature reading by the chip. Unit: degree Celsius", sensorLabels, nil)
	temperatureLabelDesc = prometheus.NewDesc(prefix+"temperature_label_info", "Suggested temperature channel label", channelLabels, nil)
	temperatureSensorEnabledDesc = prometheus.NewDesc(prefix+"temperature_sensor_enabled_bool", "1 = sensor enabled, 0 = sensor disabled", sensorLabels, nil)

	currentMaxDesc = prometheus.NewDesc(prefix+"current_max_ampere", "Current max value. Unit: Ampere", sensorLabels, nil)
	currentMinDesc = prometheus.NewDesc(prefix+"current_min_ampere", "Current min value. Unit: Ampere", sensorLabels, nil)
	currentCriticalMinValue = prometheus.NewDesc(prefix+"current_critical_min_ampere", "Current critical low value. Unit: Ampere", sensorLabels, nil)
	currentCriticalMaxValue = prometheus.NewDesc(prefix+"current_critical_max_ampere", "Current critical high value. Unit: Ampere", sensorLabels, nil)
	currentDesc = prometheus.NewDesc(prefix+"current_ampere", "Current input value. Unit: Ampere", sensorLabels, nil)
	currentSensorEnabledDesc = prometheus.NewDesc(prefix+"current_sensor_enabled_bool", "1 = sensor enabled, 0 = sensor disabled", sensorLabels, nil)

	rawValueDesc = prometheus.NewDesc(prefix+"raw_sensor_reading", "Arbitrary sensor reading, see labels on how to interpret this value", []string{"path", "description"}, nil)
}

// Collector collects hwmon metrics from the /sys filesystem
type Collector struct {
	config *Configuration
}

// NewCollector returns a new Collector instance
func NewCollector(config *Configuration) *Collector {
	return &Collector{
		config: config,
	}
}

// Describe implements collector.Collector interface Describe function
func (*Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- voltageMinDesc
	ch <- voltageCriticalMinDesc
	ch <- voltageMaxDesc
	ch <- voltageCriticalMaxDesc
	ch <- voltageDesc
	ch <- voltageLabelInfoDesc
	ch <- voltageSensorEnabledDesc

	ch <- fanMinDesc
	ch <- fanMaxDesc
	ch <- fanDesc
	ch <- fanDivisorDesc
	ch <- fanPulsesDesc
	ch <- fanTargetDesc
	ch <- fanLabelsDesc
	ch <- fanSensorEnabledDesc

	ch <- temperatureTypeDesc
	ch <- temperatureMaxDesc
	ch <- temperatureMinDesc
	ch <- temperatureMaxHysteresisDesc
	ch <- temperatureMinHysteresisDesc
	ch <- temperatureDesc
	ch <- temperatureCriticalMaxDesc
	ch <- temperatureCriticalMaxHysteresisDesc
	ch <- temperatureEmergencyMaxDesc
	ch <- temperatureEmergencyMaxHysteresisDesc
	ch <- temperatureCriticalMinDesc
	ch <- temperatureCriticalMinHysteresisDesc
	ch <- temperatureOffsetDesc
	ch <- temperatureLabelDesc
	ch <- temperatureSensorEnabledDesc

	ch <- currentMaxDesc
	ch <- currentMinDesc
	ch <- currentCriticalMinValue
	ch <- currentCriticalMaxValue
	ch <- currentDesc
	ch <- currentSensorEnabledDesc
}

type parserFunc func(string, string, string) prometheus.Metric

func getParsers(sensorType string) []parserFunc {
	return map[string][]parserFunc{
		"voltage": []parserFunc{
			makeDefaultParser(voltageMinDesc, "_min", 1000),
			makeDefaultParser(voltageCriticalMinDesc, "_lcrit", 1000),
			makeDefaultParser(voltageMaxDesc, "_max", 1000),
			makeDefaultParser(voltageDesc, "_input", 1000),
			makeChannelParser(voltageLabelInfoDesc, "_label"),
			makeDefaultParser(voltageSensorEnabledDesc, "_enable", 1),
		},
		"fan": []parserFunc{
			makeDefaultParser(fanMinDesc, "_min", 1),
			makeDefaultParser(fanMaxDesc, "_max", 1),
			makeDefaultParser(fanDesc, "_input", 1),
			makeDefaultParser(fanDivisorDesc, "_div", 1),
			makeDefaultParser(fanPulsesDesc, "_pulses", 1),
			makeDefaultParser(fanTargetDesc, "_target", 1),
			makeChannelParser(fanLabelsDesc, "_label"),
			makeDefaultParser(fanSensorEnabledDesc, "_enable", 1),
		},
		"temp": []parserFunc{
			makeTemperatureSensorTypeParser(temperatureTypeDesc, "_type"),
			makeDefaultParser(temperatureMaxDesc, "_max", 1000),
			makeDefaultParser(temperatureMinDesc, "_min", 1000),
			makeDefaultParser(temperatureMaxHysteresisDesc, "_max_hyst", 1000),
			makeDefaultParser(temperatureMinHysteresisDesc, "_min_hyst", 1000),
			makeDefaultParser(temperatureDesc, "_input", 1000),
			makeDefaultParser(temperatureCriticalMaxDesc, "_crit", 1000),
			makeDefaultParser(temperatureCriticalMaxHysteresisDesc, "_crit_hyst", 1000),
			makeDefaultParser(temperatureEmergencyMaxDesc, "_emergency", 1000),
			makeDefaultParser(temperatureEmergencyMaxHysteresisDesc, "_emergency_hyst", 1000),
			makeDefaultParser(temperatureCriticalMinDesc, "_lcrit", 1000),
			makeDefaultParser(temperatureCriticalMinHysteresisDesc, "_lcrit_hyst", 1000),
			makeDefaultParser(temperatureOffsetDesc, "_offset", 1000),
			makeChannelParser(temperatureLabelDesc, "_label"),
			makeDefaultParser(temperatureSensorEnabledDesc, "_enable", 1),
		},
		"current": []parserFunc{
			makeDefaultParser(currentMaxDesc, "_max", 1000),
			makeDefaultParser(currentMinDesc, "_min", 1000),
			makeDefaultParser(currentCriticalMinValue, "_lcrit", 1000),
			makeDefaultParser(currentCriticalMaxValue, "_crit", 1000),
			makeDefaultParser(currentDesc, "_input", 1000),
			makeDefaultParser(currentSensorEnabledDesc, "_enable", 1),
		},
		"raw": []parserFunc{
			makeRawParser(rawValueDesc),
		},
	}[sensorType]
}

func makeDefaultParser(metricDesc *prometheus.Desc, pathSuffix string, divisor float64) parserFunc {
	return func(driverPath string, hwmon string, description string) prometheus.Metric {
		value, err := util.ReadFloat64FromFile(driverPath + "/" + hwmon + pathSuffix)
		if err == nil {
			return prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value/divisor, driverPath, hwmon, description)
		}
		return nil
	}
}

func makeChannelParser(metricDesc *prometheus.Desc, pathSuffix string) parserFunc {
	return func(driverPath string, hwmon string, description string) prometheus.Metric {
		text, err := util.ReadStringFromFile(driverPath + "/" + hwmon + pathSuffix)
		if err == nil {
			return prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, 1.0, driverPath, hwmon, description, text)
		}
		return nil
	}
}

func makeTemperatureSensorTypeParser(metricDesc *prometheus.Desc, pathSuffix string) parserFunc {
	return func(driverPath string, hwmon string, description string) prometheus.Metric {
		value, err := util.ReadFloat64FromFile(driverPath + "/" + hwmon + pathSuffix)
		if err != nil {
			return nil
		}
		sensorType := map[float64]string{
			1: "CPU embedded diode",
			2: "3904 transistor",
			3: "thermal diode",
			4: "thermistor",
			5: "AMD AMDSI",
			6: "Intel PECI",
		}[value]
		return prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, 1.0, driverPath, hwmon, description, sensorType)
	}
}

func makeRawParser(metricDesc *prometheus.Desc) parserFunc {
	return func(driverPath string, hwmon string, description string) prometheus.Metric {
		value, err := util.ReadFloat64FromFile(driverPath)
		if err == nil {
			return prometheus.MustNewConstMetric(metricDesc, prometheus.GaugeValue, value, driverPath, description)
		}
		return nil
	}
}

// Name returns the string "HwmonCollector"
func (*Collector) Name() string {
	return "HwmonCollector"
}

// Collect implements collector.Collector interface Collect function
func (c *Collector) Collect(metrics chan<- prometheus.Metric, errorChan chan error, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(len(c.config.Sensors))

	for _, sensorConfiguration := range c.config.Sensors {
		collectSensor(sensorConfiguration, metrics, wg)
	}

	//wg.Wait()
}

func collectSensor(sensorConfig *SensorConfiguration, metrics chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	applicableParsers := getParsers(sensorConfig.Type)

	for _, parser := range applicableParsers {
		metric := parser(sensorConfig.DriverPath, sensorConfig.DriverHwmon, sensorConfig.Description)

		if metric != nil {
			metrics <- metric
		}
	}
}
