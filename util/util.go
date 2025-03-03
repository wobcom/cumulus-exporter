package util

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
	"strings"
)

// ReadFloat64FromFile reads a float64 from a file
func ReadFloat64FromFile(filename string) (float64, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return -1, errors.Wrapf(err, "Could not read file '%s': %v", filename, err)
	}
	valueString := strings.TrimSuffix(string(data), "\n")
	return strconv.ParseFloat(valueString, 64)
}

// ReadStringFromFile reads a sring from a file and removes a trailing newline
func ReadStringFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}
