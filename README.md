# Cumulus Exporter
[![CI Status](https://gitlab.com/wobcom/cumulus-exporter/badges/master/pipeline.svg)](https://gitlab.com/wobcom/cumulus-exporter/pipelines) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**cumulus-exporter** is a [Prometheus](https://github.com/prometheus/prometheus) exporter, that exposes metrics from switching and routing platforms running on [Cumulus Linux](https://cumulusnetworks.com/products/cumulus-linux/) based hosts alongside the [node_exporter](https://github.com/prometheus/node_exporter).

> Do not yet use the ASIC collector, as there is a known bug in current cumulus releases.

It provides the following metrics:
* Transceiver statistics (RX / TX power, voltage, temperature, ...) by including the [transceiver-exporter](https://github.com/wobcom/transceiver-exporter)
* MSTPD statistics (port (forwarding) states)
* ASIC statistics as exposed in `/cumulus/switchd`
* HWMON statistics (*needs configuration*)

## Usage
```
Usage of cumulus-exporter:
  -collecotrs.transceiver.exclude-interfaces string
    	Comma seperated list of interfaces to exclude from scrape
  -collectors.asic
    	Enable ASIC collector
  -collectors.hwmon
    	Enable hwmon collector
  -collectors.hwmon.config string
    	hwmon collector config file (default "hwmon.yml")
  -collectors.mstpd
    	Enable mstpd collector
  -collectors.mstpd.mstpctl-path string
    	mstpctl binary path (default "/sbin/mstpctl")
  -collectors.transceiver
    	Enable transceiver collector (rx / tx power, temperatures, etc.)
  -collectors.transceiver.interface-features
    	Collect interface features (results in many time series
  -version
    	Print version and exit
  -web.listen-address string
    	Address to listen on (default "[::]:9457")
  -web.telemetry-path string
    	Path under which to expose metrics (default "/metrics")
```

## Hwmon configuration
The kernel's hwmon interface provides no means of enumerating sensors, so you have to provide some information.

We have included hwmon configurations for some common models supported by Cumulus Linux in hwmon-configurations.

In case your device is not listed in here, you can figure out some sensors by running `smonctl -v --json`.

Supply the following format to the cumulus-exporter:
```yaml
sensors:
  - description: "Asic Temp Sensor"
    driver_path: "/sys/class/hwmon/hwmon8"
    driver_hwmon: temp1
    type: temp
  - description: "Main Board Ambient Sensor"
    driver_path: "/sys/class/hwmon/hwmon5"
    driver_hwmon: temp1
    type: temp
  - description: "Port Ambient Sensor"
    driver_path: "/sys/class/hwmon/hwmon4"
    driver_hwmon: temp1
    type: temp
```

If you wish to export the (floating point) content from a single file, use
```yaml
  - description: "Fan Tray 2 OK"
    driver_path: "/sys/bus/i2c/devices/0-0060/fan2_ok"
    type: raw
```
