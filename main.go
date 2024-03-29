package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/wobcom/transceiver-exporter/transceiver-collector"
	"gitlab.com/wobcom/cumulus-exporter/asic"
	"gitlab.com/wobcom/cumulus-exporter/collector"
	"gitlab.com/wobcom/cumulus-exporter/hwmon"
	"gitlab.com/wobcom/cumulus-exporter/mstpd"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version string = "1.0"

var (
	showVersion              = flag.Bool("version", false, "Print version and exit")
	listenAddress            = flag.String("web.listen-address", "[::]:9457", "Address to listen on")
	metricsPath              = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	asicCollector            = flag.Bool("collectors.asic", false, "Enable ASIC collector")
	transceiverCollector     = flag.Bool("collectors.transceiver", false, "Enable transceiver collector (rx / tx power, temperatures, etc.)")
	collectInterfaceFeatures = flag.Bool("collectors.transceiver.interface-features", false, "Collect interface features (results in many time series")
	excludeInterfaces        = flag.String("collectors.transceiver.exclude-interfaces", "", "Comma seperated list of interfaces to exclude from scrape")
	hwmonCollector           = flag.Bool("collectors.hwmon", false, "Enable hwmon collector")
	hwmonCollectorConfig     = flag.String("collectors.hwmon.config", "hwmon.yml", "hwmon collector config file")
	mstpdCollector           = flag.Bool("collectors.mstpd", false, "Enable mstpd collector")
	mstpctlPath              = flag.String("collectors.mstpd.mstpctl-path", "/sbin/mstpctl", "mstpctl binary path")
	enabledCollectors        []collector.Collector
)

func printVersion() {
	fmt.Println("cumulus-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): @fluepke")
	fmt.Println("Exposes varies metrics from devices running the Cumulus Linux operating system")
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}
	startServer()
}

func initialize() {
	if *asicCollector {
		log.Info("asic collector enabled")
		enabledCollectors = append(enabledCollectors, asic.NewCollector())
	}
	if *transceiverCollector {
		log.Info("transceiver collector enabled")
		blacklistedIfaceNames := strings.Split(*excludeInterfaces, ",")
		for index, blacklistedIfaceName := range blacklistedIfaceNames {
			blacklistedIfaceNames[index] = strings.Trim(blacklistedIfaceName, " ")
		}
		enabledCollectors = append(enabledCollectors, transceivercollector.NewCollector(blacklistedIfaceNames, *collectInterfaceFeatures, false))
	}
	if *hwmonCollector {
		log.Info("hwmon collector enabled")
		hwmonCollectorConfig, err := hwmon.LoadConfiguration(*hwmonCollectorConfig)
		if err != nil {
			log.Errorf("Could not load hwmon collector config file: %v. Disabling hwmon collector.", err)
		} else {
			enabledCollectors = append(enabledCollectors, hwmon.NewCollector(hwmonCollectorConfig))
		}
	}
	if *mstpdCollector {
		enabledCollectors = append(enabledCollectors, mstpd.NewCollector(*mstpctlPath))
	}
}

func startServer() {
	log.Infof("Starting cumulus-exporter (version: %s)", version)
	initialize()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
            <head><title>cumulus-exporter (Version ` + version + `)</title></head>
            <body>
            <h1>cumulus-exporter</h1>
            <p><a href="` + *metricsPath + `">Metrics</a></p>
            </body>
            </html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, request *http.Request) {
	registry := prometheus.NewRegistry()

	registry.MustRegister(newCumulusCollector())

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	}).ServeHTTP(w, request)
}
