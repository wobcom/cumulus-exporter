package mstpd

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"os/exec"
)

// GetBridges returns a list of the system's bridge interface names
func GetBridges() ([]string, error) {
	var res []string

	handle, err := netlink.NewHandle()
	if err != nil {
		return res, errors.Wrap(err, "Could not get netlink handle")
	} else {
		defer handle.Close()
	}

	linkList, err := handle.LinkList()
	if err != nil {
		return res, errors.Wrap(err, "Could not get system link list")
	}

	for _, link := range linkList {
		_, isBridge := link.(*netlink.Bridge)
		if !isBridge {
			continue
		}
		res = append(res, link.Attrs().Name)
	}
	return res, nil
}

// ShowPortDetail executes and parses "mstpctl showportdetails <bridge> json"
func ShowPortDetail(mstpctlPath string, bridgeName string) (ShowPortDetailResult, error) {
	cmd := exec.Command(mstpctlPath, "showportdetail", bridgeName, "json")
	res := ShowPortDetailResult{}
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer
	err := cmd.Run()
	if err != nil {
		return res, errors.Wrapf(err, "Executing '%s showportdetail %s json' failed, stderr reads: %s", mstpctlPath, bridgeName, stderrBuffer.String())
	}

	err = json.Unmarshal(stdoutBuffer.Bytes(), &res)
	if err != nil {
		return res, errors.Wrap(err, "JSON unmarshal failed")
	}

	return res, nil
}

// BoolToFloat64 returns 1 for true and 0 for false
func BoolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0
}
