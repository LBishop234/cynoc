package analysis

import (
	"main/src/domain"
)

func basicLatency(conf domain.SimConfig, tfr trafficFlowAndRoute) int {
	noFlits := tfr.PacketSize
	processingDelay := len(tfr.Route) * conf.ProcessingDelay
	return noFlits + processingDelay
}
