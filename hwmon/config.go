package hwmon

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Configuration passed to NewCollector
type Configuration struct {
	Sensors []*SensorConfiguration `yaml:"sensors"`
}

// SensorConfiguration sensor to scrape
type SensorConfiguration struct {
	Description string `yaml:"description"`
	DriverPath  string `yaml:"driver_path"`
	DriverHwmon string `yaml:"driver_hwmon,omitempty"`
	Type        string `yaml:"type"`
}

// LoadConfiguration loads and returns a configuration from a given filepath
func LoadConfiguration(path string) (*Configuration, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	configuration := &Configuration{}
	err = yaml.Unmarshal(file, configuration)
	return configuration, err
}
