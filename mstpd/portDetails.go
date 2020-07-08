package mstpd

// ShowPortDetailResult stores the parsed information of a `mstpctl showportdetail` command
type ShowPortDetailResult map[string]map[string]*PortDetails

// PortDetails store parsed port information returned for a single port by `mstpctl showportdetail`
type PortDetails struct {
	BridgeName        string  `json:"bridgeName"`
	PortName          string  `json:"portName"`
	Enabled           bool    `json:"enabled"`
	Role              string  `json:"role"`
	State             string  `json:"state"`
	ExtPortCost       float64 `json:"extPortCost"`
	AdminExtPortCost  float64 `json:"adminExtPortCost"`
	IntPortCost       float64 `json:"intPortCost"`
	AdminIntPortCost  float64 `json:"adminIntPortCost"`
	DsgnRoot          string  `json:"dsgnRoot"`
	DsgnExtCost       float64 `json:"dsgnExtCost"`
	DsgnRegRoot       string  `json:"dsgnRegRoot"`
	DsgnIntCost       float64 `json:"dsgnIntCost"`
	DsgnBr            string  `json:"dsgnBr"`
	DsgnPort          string  `json:"dsgnPort"`
	AdminEdgePort     bool    `json:"adminEdgePort,omitempty"`
	AutoEdgePort      bool    `json:"autoEdgePort"`
	OperEdgePort      bool    `json:"operEdgePort"`
	PointToPoint      bool    `json:"pointToPoint"`
	AdminPointToPoint string  `json:"adminPointToPoint"`
	PortHelloTime     float64 `json:"portHelloTime"`
	BpduGuardPort     bool    `json:"bpduGuardPort,omitempty"`
	NumTxBpdu         float64 `json:"numTxBpdu"`
	NumTxTcn          float64 `json:"numTxTcn"`
	NumRxBpdu         float64 `json:"numRxBpdu"`
	NumRxTcn          float64 `json:"numRxTcn"`
	NumTransFwd       float64 `json:"numTransFwd"`
	NumTransBlk       float64 `json:"numTransBlk"`
	BpduFilterPort    bool    `json:"bpduFilterPort,omitempty"`
	ClagRole          string  `json:"clagRole"`
	ClagDualConnMac   string  `json:"clagDualConnMac"`
	ClagRemotePortID  string  `json:"clagRemotePortId"`
	ClagSystemMac     string  `json:"clagSystemMac"`
}
