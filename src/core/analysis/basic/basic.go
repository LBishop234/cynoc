package basic

import (
	"main/src/core/analysis/util"
	"main/src/domain"
	"math"
)

func BasicLatency(conf domain.SimConfig, tfr util.TrafficFlowAndRoute) int {
	noFlits := tfr.PacketSize + util.NoAdditionalFlits

	transmission := float64(noFlits) / float64(conf.LinkBandwidth)

	headerProcessing := len(tfr.Route) * conf.ProcessingDelay

	return int(math.Round(transmission)) + headerProcessing
}
