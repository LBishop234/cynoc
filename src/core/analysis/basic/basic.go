package basic

import (
	"math"

	"main/src/core/analysis/util"
	"main/src/domain"
)

func BasicLatency(conf domain.SimConfig, tfr util.TrafficFlowAndRoute) int {
	noFlits := math.Ceil((float64(tfr.PacketSize) + float64(util.NoAdditionalFlits*conf.FlitSize)) / float64(conf.FlitSize))

	linkFlitBandwidth := float64(conf.FlitSize) / float64(conf.LinkBandwidth)

	hops := float64(len(tfr.Route))

	transmission := math.Ceil(noFlits * linkFlitBandwidth)

	headerProcessing := hops * float64(conf.ProcessingDelay)

	return int(transmission + headerProcessing)
}
