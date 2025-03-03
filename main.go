package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
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

const version string = "1.0.10"

var (
	showVersion              = flag.Bool("version", false, "Print version and exit")
	listenAddress            = flag.String("web.listen-address", "[::]:9457", "Address to listen on")
	metricsPath              = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	asicCollector            = flag.Bool("collectors.asic", false, "Enable ASIC collector")
	transceiverCollector     = flag.Bool("collectors.transceiver", false, "Enable transceiver collector (rx / tx power, temperatures, etc.)")
	collectInterfaceFeatures = flag.Bool("collectors.transceiver.interface-features", false, "Collect interface features (results in many time series")
	excludeInterfaces        = flag.String("collectors.transceiver.exclude-interfaces", "", "Comma seperated list of interfaces to exclude from scrape")
	includeInterfaces        = flag.String("collectors.transceiver.include-interfaces", "", "Comma seperated list of interfaces to include from scrape")
	excludeInterfacesRegex   = flag.String("collectors.transceiver.exclude-interfaces-regex", "", "Regex Expression for interfaces to exclude from scrape")
	includeInterfacesRegex   = flag.String("collectors.transceiver.include-interfaces-regex", "", "Regex Expression for interfaces to include from scrape")
	hwmonCollector           = flag.Bool("collectors.hwmon", false, "Enable hwmon collector")
	hwmonCollectorConfig     = flag.String("collectors.hwmon.config", "hwmon.yml", "hwmon collector config file")
	mstpdCollector           = flag.Bool("collectors.mstpd", false, "Enable mstpd collector")
	mstpctlPath              = flag.String("collectors.mstpd.mstpctl-path", "/sbin/mstpctl", "mstpctl binary path")
	logLevel                 = flag.String("log.level", "info", "The level the application logs at")
	enabledCollectors        []collector.Collector
)

func printVersion() {
	fmt.Println("cumulus-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): @fluepke, @jwagner")
	fmt.Println("Exposes varies metrics from devices running the Cumulus Linux operating system")
}

func setLogLevel() {
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		fmt.Printf("log level %s is unknown: %v\n", level, err)
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func main() {
	flag.Parse()
	setLogLevel()

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

		includedIfaceNames := strings.Split(*includeInterfaces, ",")
		for index, includedIfaceName := range includedIfaceNames {
			includedIfaceNames[index] = strings.Trim(includedIfaceName, " ")
		}

		includeIfaceRegex, includeErr := regexp.Compile(*includeInterfacesRegex)
		excludeIfaceRegex, excludeErr := regexp.Compile(*excludeInterfacesRegex)
		if includeErr != nil {
			log.Errorf("Could not compile include interface regex expression \"%s\". Disabling transceiver collector.", includeErr)
		} else if excludeErr != nil {
			log.Errorf("Could not compile exlude interface regex expression \"%s\". Disabling transceiver collector.", excludeErr)
		} else {
			enabledCollectors = append(enabledCollectors, transceivercollector.NewCollector(blacklistedIfaceNames, includedIfaceNames, includeIfaceRegex, excludeIfaceRegex, true, *collectInterfaceFeatures, false))
		}
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
