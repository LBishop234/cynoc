package basic

import (
	"main/src/core/analysis/util"
	"main/src/domain"
)

func BasicLatency(conf domain.SimConfig, tfr util.TrafficFlowAndRoute) int {
	noFlits := tfr.PacketSize
	processingDelay := len(tfr.Route) * conf.ProcessingDelay
	return noFlits + processingDelay
}
