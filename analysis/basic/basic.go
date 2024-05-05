package basic

import (
	"math"

	"main/analysis/util"
	"main/domain"
)

func BasicLatency(conf domain.SimConfig, tfr util.TrafficFlowAndRoute) int {
	noFlits := int(math.Ceil((float64(tfr.PacketSize) + float64(util.NoAdditionalFlits*conf.FlitSize)) / float64(conf.FlitSize)))

	linkBandwidth := util.LinkBandwidthFlitFactor * conf.FlitSize

	hops := len(tfr.Route)

	basicLatency := noFlits*(conf.FlitSize/linkBandwidth) + hops*conf.ProcessingDelay

	return basicLatency
}
