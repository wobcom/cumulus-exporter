package asic

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var switchdLock = &sync.Mutex{}

// ReadFloat64FromFileSwitchd reads a file from the switchd fuse
// Note: If you think this implementation is weird, this is due to
// ioutil.ReadAll crashing the switchd daemon (which would cause a downtime)
func ReadFloat64FromFileSwitchd(filename string) (float64, error) {
	switchdLock.Lock()
	time.Sleep(10 * time.Millisecond)
	defer switchdLock.Unlock()
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return -1, errors.Wrapf(err, "Could not open file %s", filename)
	}

	bytes := make([]byte, 128)
	count, err := syscall.Read(fd, bytes)
	if err != nil {
		if err1 := syscall.Close(fd); err1 != nil {
			return -1, errors.Wrapf(err, "Could not read from file and could not close file %s", filename)
		}
		return -1, errors.Wrapf(err, "Could not read from file %s", filename)
	}
	err = syscall.Close(fd)
	if err != nil {
		return -1, errors.Wrapf(err, "Could not close file %s", filename)
	}
	valueString := strings.TrimSuffix(string(bytes[:count]), "\n")
	return strconv.ParseFloat(valueString, 64)
}
