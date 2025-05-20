# Cumulus Exporter
[![CI Status](https://gitlab.com/wobcom/cumulus-exporter/badges/master/pipeline.svg)](https://gitlab.com/wobcom/cumulus-exporter/pipelines) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**cumulus-exporter** is a [Prometheus](https://github.com/prometheus/prometheus) exporter, that exposes metrics from switching and routing platforms running on [Cumulus Linux](https://cumulusnetworks.com/products/cumulus-linux/) based hosts alongside the [node_exporter](https://github.com/prometheus/node_exporter).

It provides the following metrics:
* Transceiver statistics (RX / TX power, voltage, temperature, ...) by including the [transceiver-exporter](https://github.com/wobcom/transceiver-exporter)
* MSTPD statistics (port (forwarding) states)
* ASIC statistics as exposed in `/cumulus/switchd`
* HWMON statistics (through `smonctl` utility)

## Usage
```
Usage of ./cumulus-exporter:
  -collectors.asic
    	Enable ASIC collector
  -collectors.hwmon
    	Enable hwmon collector
  -collectors.mstpd
    	Enable mstpd collector
  -collectors.mstpd.mstpctl-path string
    	mstpctl binary path (default "/sbin/mstpctl")
  -collectors.transceiver
    	Enable transceiver collector (rx / tx power, temperatures, etc.)
  -collectors.transceiver.exclude-interfaces string
    	Comma seperated list of interfaces to exclude from scrape
  -collectors.transceiver.exclude-interfaces-regex string
    	Regex Expression for interfaces to exclude from scrape
  -collectors.transceiver.include-interfaces string
    	Comma seperated list of interfaces to include from scrape
  -collectors.transceiver.include-interfaces-regex string
    	Regex Expression for interfaces to include from scrape
  -collectors.transceiver.interface-features
    	Collect interface features (results in many time series)
  -log.level string
    	The level the application logs at (default "info")
  -version
    	Print version and exit
  -web.listen-address string
    	Address to listen on (default "[::]:9457")
  -web.telemetry-path string
    	Path under which to expose metrics (default "/metrics")
```
