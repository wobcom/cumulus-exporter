package mstpd

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix = "mstpd_"

var (
	enabledDesc             *prometheus.Desc
	roleInfoDesc            *prometheus.Desc
	stateInfoDesc           *prometheus.Desc
	extPortCostDesc         *prometheus.Desc
	adminExtPortCostDesc    *prometheus.Desc
	intPortCostDesc         *prometheus.Desc
	adminIntPortCostDesc    *prometheus.Desc
	dsgnExtCostDesc         *prometheus.Desc
	dsgnIntCostDesc         *prometheus.Desc
	adminEdgePortDesc       *prometheus.Desc
	autoEdgePortDesc        *prometheus.Desc
	operEdgePortDesc        *prometheus.Desc
	pointToPointDesc        *prometheus.Desc
	portHelloTimeDesc       *prometheus.Desc
	bpduGuardPortDesc       *prometheus.Desc
	numTxBpduDesc           *prometheus.Desc
	numTxTcnDesc            *prometheus.Desc
	numRxBpduDesc           *prometheus.Desc
	numRxTcnDesc            *prometheus.Desc
	numTransFwdDesc         *prometheus.Desc
	numTransBlkDesc         *prometheus.Desc
	bpduFilterPortDesc      *prometheus.Desc
	clagRoleInfoDesc        *prometheus.Desc
	clagDualConnMacInfoDesc *prometheus.Desc
	clagSystemMacInfoDesc   *prometheus.Desc
)

// Collector collects metrics exposed by mstpctl
type Collector struct {
	mstpctlPath string
}

// NewCollector returns a new Collector instance
func NewCollector(mstpctlPath string) *Collector {
	return &Collector{
		mstpctlPath: mstpctlPath,
	}
}

func init() {
	labels := []string{"bridge_name", "interface"}
	enabledDesc = prometheus.NewDesc(prefix+"enabled_bool", "enabled", labels, nil)
	roleLabels := append(labels, "role")
	roleInfoDesc = prometheus.NewDesc(prefix+"role_info", "role", roleLabels, nil)
	stateLabels := append(labels, "state")
	stateInfoDesc = prometheus.NewDesc(prefix+"state_info", "state", stateLabels, nil)
	extPortCostDesc = prometheus.NewDesc(prefix+"ext_port_cost", "external port cost", labels, nil)
	adminExtPortCostDesc = prometheus.NewDesc(prefix+"admin_ext_port_cost", "admin external port cost", labels, nil)
	intPortCostDesc = prometheus.NewDesc(prefix+"int_port_cost", "internal port cost", labels, nil)
	adminIntPortCostDesc = prometheus.NewDesc(prefix+"admin_int_port_cost", "admin internal port cost", labels, nil)
	dsgnExtCostDesc = prometheus.NewDesc(prefix+"dsgn_ext_cost", "dsgn external cost", labels, nil)
	dsgnIntCostDesc = prometheus.NewDesc(prefix+"dsgn_int_cost", "dsgn internal cost", labels, nil)
	adminEdgePortDesc = prometheus.NewDesc(prefix+"admin_edge_port_bool", "admin edge port", labels, nil)
	autoEdgePortDesc = prometheus.NewDesc(prefix+"auto_edge_port_bool", "auto edge port", labels, nil)
	operEdgePortDesc = prometheus.NewDesc(prefix+"oper_edge_port_bool", "oper edge port", labels, nil)
	pointToPointDesc = prometheus.NewDesc(prefix+"point_to_point_bool", "point-to-point", labels, nil)
	portHelloTimeDesc = prometheus.NewDesc(prefix+"port_hello_time_seconds", "port hello time in seconds", labels, nil)
	bpduGuardPortDesc = prometheus.NewDesc(prefix+"bpdu_guard_port_bool", "bpdu guard port", labels, nil)
	numTxBpduDesc = prometheus.NewDesc(prefix+"num_tx_bpdu_total", "Num TX BPDU", labels, nil)
	numTxTcnDesc = prometheus.NewDesc(prefix+"num_tx_tcn_total", "Num TX TCN", labels, nil)
	numRxBpduDesc = prometheus.NewDesc(prefix+"num_rx_bpdu_total", "Num RX BPDU", labels, nil)
	numRxTcnDesc = prometheus.NewDesc(prefix+"num_rx_tcn_total", "Num RX TCN", labels, nil)
	numTransFwdDesc = prometheus.NewDesc(prefix+"num_trans_fw_total", "Num Transition FWD", labels, nil)
	numTransBlkDesc = prometheus.NewDesc(prefix+"num_trans_blk_total", "Num Transition BLK", labels, nil)
	bpduFilterPortDesc = prometheus.NewDesc(prefix+"bpdu_filter_port_bool", "bpdufilter port", labels, nil)
	clagRoleLabels := append(labels, "clag_role")
	clagRoleInfoDesc = prometheus.NewDesc(prefix+"clag_role_info", "clag role", clagRoleLabels, nil)
	clagDualConnMacLabels := append(labels, "clag_dual_conn_mac")
	clagDualConnMacInfoDesc = prometheus.NewDesc(prefix+"clag_dual_conn_mac_info", "clag dual conn mac", clagDualConnMacLabels, nil)
	clagSystemMacInfoLabels := append(labels, "clag_system_mac")
	clagSystemMacInfoDesc = prometheus.NewDesc(prefix+"clag_system_mac_info", "clag system mac", clagSystemMacInfoLabels, nil)
}

// Describe implements collector.Collector interface's Describe function
func (*Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- enabledDesc
	ch <- roleInfoDesc
	ch <- stateInfoDesc
	ch <- extPortCostDesc
	ch <- adminExtPortCostDesc
	ch <- intPortCostDesc
	ch <- adminIntPortCostDesc
	ch <- dsgnExtCostDesc
	ch <- dsgnIntCostDesc
	ch <- adminEdgePortDesc
	ch <- autoEdgePortDesc
	ch <- operEdgePortDesc
	ch <- pointToPointDesc
	ch <- portHelloTimeDesc
	ch <- bpduGuardPortDesc
	ch <- numTxBpduDesc
	ch <- numTxTcnDesc
	ch <- numRxBpduDesc
	ch <- numRxTcnDesc
	ch <- numTransFwdDesc
	ch <- numTransBlkDesc
	ch <- bpduFilterPortDesc
	ch <- clagRoleInfoDesc
	ch <- clagDualConnMacInfoDesc
	ch <- clagSystemMacInfoDesc
}

// Collect implements collector.Collector interface's Collect function
func (c *Collector) Collect(metrics chan<- prometheus.Metric, errorChan chan error, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	bridges, err := GetBridges()
	if err != nil {
		errorChan <- errors.Wrap(err, "Could not retrieve list of system's bridge interfaces")
	}

	for _, bridge := range bridges {
		showPortDetail, err := ShowPortDetail(c.mstpctlPath, bridge)
		if err != nil {
			errorChan <- errors.Wrapf(err, "Show port failed for interface %s", bridge)
		}
		for _, portData := range showPortDetail {
			for _, portDetails := range portData {
				collectForPort(portDetails, metrics)
			}
		}
	}
}

// Name returns the string "MstpdCollector"
func (*Collector) Name() string {
	return "MstpdCollector"
}

func collectForPort(portDetails *PortDetails, metrics chan<- prometheus.Metric) {
	labels := []string{portDetails.BridgeName, portDetails.PortName}
	metrics <- prometheus.MustNewConstMetric(enabledDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.Enabled), labels...)
	roleLabels := append(labels, portDetails.Role)
	metrics <- prometheus.MustNewConstMetric(roleInfoDesc, prometheus.GaugeValue, 1.0, roleLabels...)
	stateLabels := append(labels, portDetails.State)
	metrics <- prometheus.MustNewConstMetric(stateInfoDesc, prometheus.GaugeValue, 1.0, stateLabels...)
	metrics <- prometheus.MustNewConstMetric(extPortCostDesc, prometheus.GaugeValue, portDetails.ExtPortCost, labels...)
	metrics <- prometheus.MustNewConstMetric(adminExtPortCostDesc, prometheus.GaugeValue, portDetails.AdminExtPortCost, labels...)
	metrics <- prometheus.MustNewConstMetric(intPortCostDesc, prometheus.GaugeValue, portDetails.IntPortCost, labels...)
	metrics <- prometheus.MustNewConstMetric(adminIntPortCostDesc, prometheus.GaugeValue, portDetails.AdminIntPortCost, labels...)
	metrics <- prometheus.MustNewConstMetric(dsgnExtCostDesc, prometheus.GaugeValue, portDetails.DsgnExtCost, labels...)
	metrics <- prometheus.MustNewConstMetric(dsgnIntCostDesc, prometheus.GaugeValue, portDetails.DsgnIntCost, labels...)
	metrics <- prometheus.MustNewConstMetric(adminEdgePortDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.AdminEdgePort), labels...)
	metrics <- prometheus.MustNewConstMetric(autoEdgePortDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.AutoEdgePort), labels...)
	metrics <- prometheus.MustNewConstMetric(operEdgePortDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.OperEdgePort), labels...)
	metrics <- prometheus.MustNewConstMetric(pointToPointDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.PointToPoint), labels...)
	metrics <- prometheus.MustNewConstMetric(portHelloTimeDesc, prometheus.GaugeValue, portDetails.PortHelloTime, labels...)
	metrics <- prometheus.MustNewConstMetric(bpduGuardPortDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.BpduGuardPort), labels...)
	metrics <- prometheus.MustNewConstMetric(numTxBpduDesc, prometheus.GaugeValue, portDetails.NumTxBpdu, labels...)
	metrics <- prometheus.MustNewConstMetric(numTxTcnDesc, prometheus.GaugeValue, portDetails.NumTxTcn, labels...)
	metrics <- prometheus.MustNewConstMetric(numRxBpduDesc, prometheus.GaugeValue, portDetails.NumRxBpdu, labels...)
	metrics <- prometheus.MustNewConstMetric(numRxTcnDesc, prometheus.GaugeValue, portDetails.NumRxTcn, labels...)
	metrics <- prometheus.MustNewConstMetric(numTransFwdDesc, prometheus.GaugeValue, portDetails.NumTransFwd, labels...)
	metrics <- prometheus.MustNewConstMetric(numTransBlkDesc, prometheus.GaugeValue, portDetails.NumTransBlk, labels...)
	metrics <- prometheus.MustNewConstMetric(bpduFilterPortDesc, prometheus.GaugeValue, BoolToFloat64(portDetails.BpduFilterPort), labels...)
	clagRoleLabels := append(labels, portDetails.ClagRole)
	metrics <- prometheus.MustNewConstMetric(clagRoleInfoDesc, prometheus.GaugeValue, 1.0, clagRoleLabels...)
	clagDualConnMacLabels := append(labels, portDetails.ClagDualConnMac)
	metrics <- prometheus.MustNewConstMetric(clagDualConnMacInfoDesc, prometheus.GaugeValue, 1.0, clagDualConnMacLabels...)
	clagSystemMacInfoLabels := append(labels, portDetails.ClagSystemMac)
	metrics <- prometheus.MustNewConstMetric(clagSystemMacInfoDesc, prometheus.GaugeValue, 1.0, clagSystemMacInfoLabels...)
}
