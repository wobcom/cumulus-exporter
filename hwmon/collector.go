package hwmon

import (
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
)

const prefix = "hwmon_"

var (
	voltageMinDesc           *prometheus.Desc
	voltageCriticalMinDesc   *prometheus.Desc
	voltageMaxDesc           *prometheus.Desc
	voltageCriticalMaxDesc   *prometheus.Desc
	voltageDesc              *prometheus.Desc

	fanMinDesc           *prometheus.Desc
	fanMaxDesc           *prometheus.Desc
	fanDesc              *prometheus.Desc
	fanPulsesDesc        *prometheus.Desc
	fanTargetDesc        *prometheus.Desc

	temperatureMinDesc                    *prometheus.Desc
	temperatureMaxDesc                    *prometheus.Desc
	temperatureMinHysteresisDesc          *prometheus.Desc
	temperatureMaxHysteresisDesc          *prometheus.Desc
	temperatureDesc                       *prometheus.Desc
	temperatureCriticalMinDesc            *prometheus.Desc
	temperatureCriticalMaxDesc            *prometheus.Desc
	temperatureCriticalMinHysteresisDesc  *prometheus.Desc
	temperatureCriticalMaxHysteresisDesc  *prometheus.Desc
	temperatureEmergencyMaxDesc           *prometheus.Desc
	temperatureEmergencyMaxHysteresisDesc *prometheus.Desc
	temperatureOffsetDesc                 *prometheus.Desc

	currentMinDesc           *prometheus.Desc
	currentMaxDesc           *prometheus.Desc
	currentCriticalMinValue  *prometheus.Desc
	currentCriticalMaxValue  *prometheus.Desc
	currentDesc              *prometheus.Desc

	powerWatt      *prometheus.Desc
	powerPresent   *prometheus.Desc
	powerAllOk     *prometheus.Desc
	powerAllOkPrev *prometheus.Desc

	rawValueDesc *prometheus.Desc
)

func init() {
	sensorLabels := []string{"hw_mon", "description"}

	voltageMinDesc           = prometheus.NewDesc(prefix+"voltage_min_volts", "Voltage min value. Unit: Volts", sensorLabels, nil)
	voltageCriticalMinDesc   = prometheus.NewDesc(prefix+"voltage_critical_min_volts", "Voltage critical min value. Unit: Volts", sensorLabels, nil)
	voltageMaxDesc           = prometheus.NewDesc(prefix+"voltage_max_volts", "Voltage max value. Unit: Volts", sensorLabels, nil)
	voltageCriticalMaxDesc   = prometheus.NewDesc(prefix+"voltage_critical_max_volts", "Voltage critical max value. Unit: Volts", sensorLabels, nil)
	voltageDesc              = prometheus.NewDesc(prefix+"voltage_volts", "Voltage input value. Unit: Volts", sensorLabels, nil)

	fanMinDesc           = prometheus.NewDesc(prefix+"fan_min_rpm", "Fan minimum value. Unit: revolution/min", sensorLabels, nil)
	fanMaxDesc           = prometheus.NewDesc(prefix+"fan_max_rpm", "Fan maximum value. Unit: revolution/min", sensorLabels, nil)
	fanDesc              = prometheus.NewDesc(prefix+"fan_rpm", "Fan input value. Unit: revolution/min", sensorLabels, nil)
	fanPulsesDesc        = prometheus.NewDesc(prefix+"fan_pulses", "Number of tachometer pulses per fan revolution", sensorLabels, nil)
	fanTargetDesc        = prometheus.NewDesc(prefix+"fan_target_rpm", "Desired fan speed. Unit: revolution/min", sensorLabels, nil)

	temperatureMaxDesc                    = prometheus.NewDesc(prefix+"temperature_max_celsius", "Temperature max value. Unit: degree Celsius", sensorLabels, nil)
	temperatureMinDesc                    = prometheus.NewDesc(prefix+"temperature_min_celsius", "Temperature min value. Unit: degree Celsius", sensorLabels, nil)
	temperatureMaxHysteresisDesc          = prometheus.NewDesc(prefix+"temperature_max_hysteresis_celsius", "Temperature hysteresis value for max limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureMinHysteresisDesc          = prometheus.NewDesc(prefix+"temperature_min_hysteresis_celsius", "Temperature hysteresis value for min limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureDesc                       = prometheus.NewDesc(prefix+"temperature_celsius", "Temperature input value. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMaxDesc            = prometheus.NewDesc(prefix+"temperature_critical_max_celsius", "Temperature critical max value, typically greater than corresponding temp_max values. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMaxHysteresisDesc  = prometheus.NewDesc(prefix+"temperature_critical_max_hysteresis_celsius", "Temperature hysteresis value for critical limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureEmergencyMaxDesc           = prometheus.NewDesc(prefix+"temperature_emergency_max_celsius", "Temperature emergency max value, for chips supporting more than two upper temperature limits. Unit: degree Celsius", sensorLabels, nil)
	temperatureEmergencyMaxHysteresisDesc = prometheus.NewDesc(prefix+"temperature_emergency_max_hysteresis_celsius", "Temperature hysteresis value for emergency limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMinDesc            = prometheus.NewDesc(prefix+"temperature_critical_min_celsius", "Temperature criticial min value, typically lower than corresponding temp_min values. Unit: degree Celsius", sensorLabels, nil)
	temperatureCriticalMinHysteresisDesc  = prometheus.NewDesc(prefix+"temperature_critical_min_hysteresis_celsius", "Temperature hysteresis value for critical min limit. Unit: degree Celsius", sensorLabels, nil)
	temperatureOffsetDesc                 = prometheus.NewDesc(prefix+"temperature_offset_celsius", "Temperature offset which is added to the temperature reading by the chip. Unit: degree Celsius", sensorLabels, nil)

	currentMaxDesc           = prometheus.NewDesc(prefix+"current_max_ampere", "Current max value. Unit: Ampere", sensorLabels, nil)
	currentMinDesc           = prometheus.NewDesc(prefix+"current_min_ampere", "Current min value. Unit: Ampere", sensorLabels, nil)
	currentCriticalMinValue  = prometheus.NewDesc(prefix+"current_critical_min_ampere", "Current critical low value. Unit: Ampere", sensorLabels, nil)
	currentCriticalMaxValue  = prometheus.NewDesc(prefix+"current_critical_max_ampere", "Current critical high value. Unit: Ampere", sensorLabels, nil)
	currentDesc              = prometheus.NewDesc(prefix+"current_ampere", "Current input value. Unit: Ampere", sensorLabels, nil)

	powerWatt        = prometheus.NewDesc(prefix+"power_watt", "Current Usage. Unit: Watt ", sensorLabels, nil)
	powerPresent     = prometheus.NewDesc(prefix+"power_present", "Is Power Present. 1 = present, 0 = missing", sensorLabels, nil)
	powerAllOk       = prometheus.NewDesc(prefix+"power_all_ok", "Is PSU Ok. 1 = OK, 0 = BAD, -1 = POWERED OFF, -2 NOT DETECTED", sensorLabels, nil)
	powerAllOkPrev   = prometheus.NewDesc(prefix+"power_all_ok_prev", "Is PSU Ok (Previous State). 1 = OK, 0 = BAD, -1 = POWERED OFF, -2 NOT DETECTED", sensorLabels, nil)

	rawValueDesc = prometheus.NewDesc(prefix+"raw_sensor_reading", "Arbitrary sensor reading, see labels on how to interpret this value", []string{"description"}, nil)
}

// Collector collects hwmon metrics from the /sys filesystem
type Collector struct {}

// NewCollector returns a new Collector instance
func NewCollector() *Collector {
	return &Collector {};
}

// Describe implements collector.Collector interface Describe function
func (*Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- voltageMinDesc
	ch <- voltageCriticalMinDesc
	ch <- voltageMaxDesc
	ch <- voltageCriticalMaxDesc
	ch <- voltageDesc

	ch <- fanMinDesc
	ch <- fanMaxDesc
	ch <- fanDesc
	ch <- fanPulsesDesc
	ch <- fanTargetDesc

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

	ch <- currentMaxDesc
	ch <- currentMinDesc
	ch <- currentCriticalMinValue
	ch <- currentCriticalMaxValue
	ch <- currentDesc

	ch <- powerWatt
	ch <- powerPresent
	ch <- powerAllOk
	ch <- powerAllOkPrev
}

func (*Collector) Name() string {
	return "HwmonCollector"
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric, errorChan chan error, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	smonCtlOut, err := runSmonCtl()
	if err != nil {
		errorChan <- err
		return
	}

	collectSensors(smonCtlOut, metrics, errorChan)
}

func runSmonCtl() ([]byte, error) {
	cmd := exec.Command("smonctl", "--json", "-v")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err	
	}

	return out, nil
}

func collectSensors(data []byte, metrics chan<- prometheus.Metric, errorChan chan error) {
	sensors, err := UnmarshalSensors(data)
	if err != nil {
		errorChan <- err
	}

	for _, sensor := range sensors {
		sensor.Collect(metrics)
	}
}

