package asic

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gitlab.com/wobcom/cumulus-exporter/collector"
)

const prefix string = "cumulus_switchd_"
const statPath string = "/cumulus/switchd/run/"

var (
	host0EntriesDesc                 *prometheus.Desc
	host1EntriesDesc                 *prometheus.Desc
	route0EntriesDesc                *prometheus.Desc
	route1EntriesDesc                *prometheus.Desc
	ipv4HostEntriesDesc              *prometheus.Desc
	ipv6HostEntriesDesc              *prometheus.Desc
	ipv46HostEntriesDesc             *prometheus.Desc
	ipv4RouteEntriesDesc             *prometheus.Desc
	longIpv6RouteEntriesDesc         *prometheus.Desc
	ipv6RouteEntriesDesc             *prometheus.Desc
	ipv46RouteEntriesDesc            *prometheus.Desc
	ipv4NeighborsDesc                *prometheus.Desc
	ipv6NeighborsDesc                *prometheus.Desc
	routesTotalDesc                  *prometheus.Desc
	ecmpNexthopsDesc                 *prometheus.Desc
	macEntriesDesc                   *prometheus.Desc
	mcastRoutesTotalDesc             *prometheus.Desc
	ingressACLEntriesDesc            *prometheus.Desc
	ingressACLCountersDesc           *prometheus.Desc
	ingressACLMetersDesc             *prometheus.Desc
	ingressACLSlicesDesc             *prometheus.Desc
	egressACLEntriesDesc             *prometheus.Desc
	egressACLCountersDesc            *prometheus.Desc
	egressACLMetersDesc              *prometheus.Desc
	egressACLSlicesDesc              *prometheus.Desc
	ingressACLIPv4MACFilterTableDesc *prometheus.Desc
	ingressACLIPv6FilterTableDesc    *prometheus.Desc
	ingressACLMirrorTableDesc        *prometheus.Desc
	ingressACL8021xFilterTableDesc   *prometheus.Desc
	ingressPBRIPv4FilterTableDesc    *prometheus.Desc
	ingressPBRIPv6FilterTableDesc    *prometheus.Desc
	ingressACLIPv4MangleTableDesc    *prometheus.Desc
	ingressACLIPv6MangleTableDesc    *prometheus.Desc
	egressACLIPv4MACFilterTableDesc  *prometheus.Desc
	egressACLIPv6FilterTableDesc     *prometheus.Desc
	aclL4PortRangeCheckersDesc       *prometheus.Desc
)

type statsCallbackFunc func(metrics chan<- prometheus.Metric) error

func makeDefaultStatsCallback(metricsDesc *prometheus.Desc, countFile string, maxFile string) statsCallbackFunc {
	if maxFile == "" {
		return func(metrics chan<- prometheus.Metric) error {
			currentValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, countFile))
			if err != nil {
				return errors.Wrapf(err, "Could not read current value from file '%s': %v", countFile, err)
			}

			metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, currentValue, "current")
			return nil
		}
	}
	return func(metrics chan<- prometheus.Metric) error {
		currentValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, countFile))
		if err != nil {
			return errors.Wrapf(err, "Could not read current value from file '%s': %v", countFile, err)
		}
		maxValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, maxFile))
		if err != nil {
			return errors.Wrapf(err, "Could not read max value from file '%s': %v", maxFile, err)
		}

		metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, currentValue, "current")
		metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, maxValue, "max")
		return nil
	}
}

func makeDefaultStatsCallbackAllocations(metricsDesc *prometheus.Desc, countFile string, maxFile string, allocationsFile string) statsCallbackFunc {
	return func(metrics chan<- prometheus.Metric) error {
		currentValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, countFile))
		if err != nil {
			return errors.Wrapf(err, "Could not read current value from file '%s': %v", countFile, err)
		}
		maxValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, maxFile))
		if err != nil {
			return errors.Wrapf(err, "Could not read max value from file '%s': %v", maxFile, err)
		}
		allocatedValue, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, allocationsFile))
		if err != nil {
			return errors.Wrapf(err, "Could not read allocated value from file '%s': %v", allocationsFile, err)
		}
		metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, currentValue, "current")
		metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, maxValue, "max")
		metrics <- prometheus.MustNewConstMetric(metricsDesc, prometheus.GaugeValue, allocatedValue, "allocated")
		return nil
	}
}

func makeHost0Callback(hostMode int) statsCallbackFunc {
	if hostMode == 1 {
		return makeDefaultStatsCallback(ipv46HostEntriesDesc, "route_info/host/count_0", "route_info/host/max_0")
	} else if hostMode == 2 {
		return makeDefaultStatsCallback(ipv4HostEntriesDesc, "route_info/host/count_0", "route_info/host/max_0")
	} else {
		return makeDefaultStatsCallback(host0EntriesDesc, "route_info/host/count_0", "route_info/host/max_0")
	}
}

func makeHost1Callback(hostMode int) statsCallbackFunc {
	if hostMode == 1 {
		return func(metrics chan<- prometheus.Metric) error {
			return nil
		}
	} else if hostMode == 2 {
		return makeDefaultStatsCallback(ipv6HostEntriesDesc, "route_info/host/count_1", "route_info/host/max_1")
	} else {
		return makeDefaultStatsCallback(host1EntriesDesc, "route_info/host/count_1", "route_info/host/max_1")
	}
}

func makeRoute0Callback(routeMode int) statsCallbackFunc {
	if routeMode == 1 {
		return makeDefaultStatsCallback(ipv46RouteEntriesDesc, "route_info/route/count_0", "route_info/route/max_0")
	} else if routeMode == 2 {
		return makeDefaultStatsCallback(ipv4RouteEntriesDesc, "route_info/route/count_0", "route_info/route/max_0")
	} else {
		return makeDefaultStatsCallback(route0EntriesDesc, "route_info/route/count_0", "route_info/route/max_0")
	}
}

func makeRoute1Callback(routeMode int) statsCallbackFunc {
	if routeMode == 1 {
		return makeDefaultStatsCallback(longIpv6RouteEntriesDesc, "route_info/route/count_1", "route_info/route/max_1")
	} else if routeMode == 2 {
		return makeDefaultStatsCallback(ipv6RouteEntriesDesc, "route_info/route/count_1", "route_info/route/max_1")
	} else {
		return makeDefaultStatsCallback(route1EntriesDesc, "route_info/route/count_1", "route_info/route/max_1")
	}
}

func makeStats(routeMode int, hostMode int) []statsCallbackFunc {
	return []statsCallbackFunc{
		makeHost0Callback(hostMode),
		makeHost1Callback(hostMode),
		makeRoute0Callback(routeMode),
		makeRoute1Callback(routeMode),
		makeDefaultStatsCallback(ipv4NeighborsDesc, "route_info/host/count_v4", ""),
		makeDefaultStatsCallback(ipv6NeighborsDesc, "route_info/host/count_v6", ""),
		makeDefaultStatsCallback(routesTotalDesc, "route_info/route/count_total", "route_info/route/max_total"),
		makeDefaultStatsCallback(ecmpNexthopsDesc, "route_info/ecmp_nh/count", "route_info/ecmp_nh/max"),
		makeDefaultStatsCallback(macEntriesDesc, "route_info/mac/count", "route_info/mac/max"),
		makeDefaultStatsCallback(mcastRoutesTotalDesc, "route_info/mroute/count_total", "route_info/mroute/max_total"),
		makeDefaultStatsCallback(ingressACLEntriesDesc, "acl_info/ingress/entries", "acl_info/ingress/entries_total"),
		makeDefaultStatsCallback(ingressACLCountersDesc, "acl_info/ingress/counters", "acl_info/ingress/counters_total"),
		makeDefaultStatsCallback(ingressACLMetersDesc, "acl_info/ingress/meters", "acl_info/ingress/meters_total"),
		makeDefaultStatsCallback(ingressACLSlicesDesc, "acl_info/ingress/slices", "acl_info/ingress/slices_total"),
		makeDefaultStatsCallback(egressACLEntriesDesc, "acl_info/egress/entries", "acl_info/egress/entries_total"),
		makeDefaultStatsCallback(egressACLCountersDesc, "acl_info/egress/counters", "acl_info/egress/counters_total"),
		makeDefaultStatsCallback(egressACLMetersDesc, "acl_info/egress/meters", "acl_info/egress/meters_total"),
		makeDefaultStatsCallback(egressACLSlicesDesc, "acl_info/egress/slices", "acl_info/egress/slices_total"),
		makeDefaultStatsCallbackAllocations(ingressACLIPv4MACFilterTableDesc, "acl_info/ingress/v4mac_filter/entries_used", "acl_info/ingress/v4mac_filter/entries_max", "acl_info/ingress/v4mac_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressACLIPv6FilterTableDesc, "acl_info/ingress/v6_filter/entries_used", "acl_info/ingress/v6_filter/entries_max", "acl_info/ingress/v6_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressACLMirrorTableDesc, "acl_info/ingress/mirror_filter/entries_used", "acl_info/ingress/mirror_filter/entries_max", "acl_info/ingress/mirror_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressACL8021xFilterTableDesc, "acl_info/ingress/mirror_filter/entries_used", "acl_info/ingress/mirror_filter/entries_max", "acl_info/ingress/mirror_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressPBRIPv4FilterTableDesc, "iprule/info/ingress/v4mac_filter/entries_used", "iprule/info/ingress/v4mac_filter/entries_max", "iprule/info/ingress/v4mac_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressPBRIPv6FilterTableDesc, "iprule/info/ingress/v6_filter/entries_used", "iprule/info/ingress/v6_filter/entries_max", "iprule/info/ingress/v6_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressACLIPv4MangleTableDesc, "acl_info/ingress/v4mac_mangle/entries_used", "acl_info/ingress/v4mac_mangle/entries_max", "acl_info/ingress/v4mac_mangle/entries_allocated"),
		makeDefaultStatsCallbackAllocations(ingressACLIPv6MangleTableDesc, "acl_info/ingress/v6_mangle/entries_used", "acl_info/ingress/v6_mangle/entries_max", "acl_info/ingress/v6_mangle/entries_allocated"),
		makeDefaultStatsCallbackAllocations(egressACLIPv4MACFilterTableDesc, "acl_info/egress/v4mac_filter/entries_used", "acl_info/egress/v4mac_filter/entries_max", "acl_info/egress/v4mac_filter/entries_allocated"),
		makeDefaultStatsCallbackAllocations(egressACLIPv6FilterTableDesc, "acl_info/egress/v6_filter/entries_used", "acl_info/egress/v6_filter/entries_max", "acl_info/egress/v6_filter/entries_allocated"),
		makeDefaultStatsCallback(aclL4PortRangeCheckersDesc, "acl_info/l4_port_range_checkers/entries_used", "acl_info/l4_port_range_checkers/entries_max"),
	}
}

// Collector collects metrics exposed by switchd in the /cumulus/switchd fuse
type Collector struct{}

// NewCollector returns a new Collector instance
func NewCollector() collector.Collector {
	return &Collector{}
}

// Name returns the string "AsicCollector"
func (*Collector) Name() string {
	return "AsicCollector"
}

func init() {
	labels := []string{"reading_type"}
	host0EntriesDesc = prometheus.NewDesc(prefix+"host_0_entry", "Host 0 entries", labels, nil)
	host1EntriesDesc = prometheus.NewDesc(prefix+"host_1_entry", "Host 0 entries", labels, nil)
	route0EntriesDesc = prometheus.NewDesc(prefix+"route_0_entry", "Route 0 entries", labels, nil)
	route1EntriesDesc = prometheus.NewDesc(prefix+"route_1_entry", "Route 1 entries", labels, nil)
	ipv4HostEntriesDesc = prometheus.NewDesc(prefix+"host_v4_entry", "IPv4 host entries", labels, nil)
	ipv6HostEntriesDesc = prometheus.NewDesc(prefix+"host_v6_entry", "IPv6 host entries", labels, nil)
	ipv46HostEntriesDesc = prometheus.NewDesc(prefix+"host_v46_entry", "IPv4/IPv6 host entries", labels, nil)
	ipv4RouteEntriesDesc = prometheus.NewDesc(prefix+"route_v4_entry", "IPv4 route entries", labels, nil)
	longIpv6RouteEntriesDesc = prometheus.NewDesc(prefix+"long_route_v6_entry", "Long IPv6 route entries", labels, nil)
	ipv6RouteEntriesDesc = prometheus.NewDesc(prefix+"route_v6_entry", "IPv6 route entries", labels, nil)
	ipv46RouteEntriesDesc = prometheus.NewDesc(prefix+"route_v46_entry", "IPv4/IPv6 route entries", labels, nil)
	ipv4NeighborsDesc = prometheus.NewDesc(prefix+"neighbor_v4_entry", "IPv4 neighbors", labels, nil)
	ipv6NeighborsDesc = prometheus.NewDesc(prefix+"neighbor_v6_entry", "IPv6 neighbors", labels, nil)
	routesTotalDesc = prometheus.NewDesc(prefix+"route_total_entry", "Total Routes", labels, nil)
	ecmpNexthopsDesc = prometheus.NewDesc(prefix+"ecmp_nh_entry", "ECMP nexthops", labels, nil)
	macEntriesDesc = prometheus.NewDesc(prefix+"mac_entry", "MAC entries", labels, nil)
	mcastRoutesTotalDesc = prometheus.NewDesc(prefix+"mroute_total_entry", "Total Mcast Routes", labels, nil)
	ingressACLEntriesDesc = prometheus.NewDesc(prefix+"in_acl_entry", "Ingress ACL entries", labels, nil)
	ingressACLCountersDesc = prometheus.NewDesc(prefix+"in_acl_counter", "Ingress ACL counters", labels, nil)
	ingressACLMetersDesc = prometheus.NewDesc(prefix+"in_acl_meter", "Ingress ACL meters", labels, nil)
	ingressACLSlicesDesc = prometheus.NewDesc(prefix+"in_acl_slice", "Ingress ACL slices", labels, nil)
	egressACLEntriesDesc = prometheus.NewDesc(prefix+"eg_acl_entry", "Egress ACL entries", labels, nil)
	egressACLCountersDesc = prometheus.NewDesc(prefix+"eg_acl_counter", "Egress ACL counters", labels, nil)
	egressACLMetersDesc = prometheus.NewDesc(prefix+"eg_acl_meter", "Egress ACL meters", labels, nil)
	egressACLSlicesDesc = prometheus.NewDesc(prefix+"eg_acl_slice", "Egress ACL slices", labels, nil)
	ingressACLIPv4MACFilterTableDesc = prometheus.NewDesc(prefix+"in_acl_v4mac_filter", "Ingress ACL ipv4_mac filter table", labels, nil)
	ingressACLIPv6FilterTableDesc = prometheus.NewDesc(prefix+"in_acl_v6_filter", "Ingress ACL ipv6 filter table", labels, nil)
	ingressACLMirrorTableDesc = prometheus.NewDesc(prefix+"in_acl_mirror_filter", "Ingress ACL mirror table", labels, nil)
	ingressACL8021xFilterTableDesc = prometheus.NewDesc(prefix+"in_acl_8021x_filter", "Ingress ACL 8021x filter table", labels, nil)
	ingressPBRIPv4FilterTableDesc = prometheus.NewDesc(prefix+"in_pbr_v4mac_filter", "Ingress PBR ipv4_mac filter table", labels, nil)
	ingressPBRIPv6FilterTableDesc = prometheus.NewDesc(prefix+"in_pbr_v6_filter", "Ingress PBR ipv6 filter table", labels, nil)
	ingressACLIPv4MangleTableDesc = prometheus.NewDesc(prefix+"in_acl_v4mac_mangle", "Ingress ACL ipv4_mac mangle table", labels, nil)
	ingressACLIPv6MangleTableDesc = prometheus.NewDesc(prefix+"in_acl_v6_mangle", "Ingress ACL ipv6 mangle table", labels, nil)
	egressACLIPv4MACFilterTableDesc = prometheus.NewDesc(prefix+"eg_acl_v4mac_filter", "Egress ACL ipv4_mac filter table", labels, nil)
	egressACLIPv6FilterTableDesc = prometheus.NewDesc(prefix+"eg_acl_v6_filter", "Egress ACL ipv6 filter table", labels, nil)
	aclL4PortRangeCheckersDesc = prometheus.NewDesc(prefix+"acl_l4_port_range_checkers", "ACL L4 port range checkers", labels, nil)
}

// Describe implements collector.Collector interface's Describe function
func (*Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- host0EntriesDesc
	ch <- host1EntriesDesc
	ch <- route0EntriesDesc
	ch <- route1EntriesDesc
	ch <- ipv4HostEntriesDesc
	ch <- ipv6HostEntriesDesc
	ch <- ipv46HostEntriesDesc
	ch <- ipv4RouteEntriesDesc
	ch <- longIpv6RouteEntriesDesc
	ch <- ipv6RouteEntriesDesc
	ch <- ipv46RouteEntriesDesc
	ch <- ipv4NeighborsDesc
	ch <- ipv6NeighborsDesc
	ch <- routesTotalDesc
	ch <- ecmpNexthopsDesc
	ch <- macEntriesDesc
	ch <- mcastRoutesTotalDesc
	ch <- ingressACLEntriesDesc
	ch <- ingressACLCountersDesc
	ch <- ingressACLMetersDesc
	ch <- ingressACLSlicesDesc
	ch <- egressACLEntriesDesc
	ch <- egressACLCountersDesc
	ch <- egressACLMetersDesc
	ch <- egressACLSlicesDesc
	ch <- ingressACLIPv4MACFilterTableDesc
	ch <- ingressACLIPv6FilterTableDesc
	ch <- ingressACLMirrorTableDesc
	ch <- ingressACL8021xFilterTableDesc
	ch <- ingressPBRIPv4FilterTableDesc
	ch <- ingressPBRIPv6FilterTableDesc
	ch <- ingressACLIPv4MangleTableDesc
	ch <- ingressACLIPv6MangleTableDesc
	ch <- egressACLIPv4MACFilterTableDesc
	ch <- egressACLIPv6FilterTableDesc
	ch <- aclL4PortRangeCheckersDesc
}

func getRouteMode() (int, error) {
	routeMode, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, "route_info/route/mode"))
	return int(routeMode), err
}

func getHostMode() (int, error) {
	hostMode, err := ReadFloat64FromFileSwitchd(filepath.Join(statPath, "route_info/host/mode"))
	return int(hostMode), err
}

// Collect implements collector.Collector interface's Collect function
func (c *Collector) Collect(metrics chan<- prometheus.Metric, errorChan chan error, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	routeMode, err := getRouteMode()
	if err != nil {
		errorChan <- errors.Wrapf(err, "Could not retrieve route mode: %v", err)
	}

	hostMode, err := getHostMode()
	if err != nil {
		errorChan <- errors.Wrapf(err, "Could not retrieve host mode: %v", err)
	}

	stats := makeStats(routeMode, hostMode)

	log.Infof("len(stats) = %d", len(stats))
	for _, stat := range stats {
		err := stat(metrics)
		if err != nil {
			errorChan <- err
		}
	}
}
