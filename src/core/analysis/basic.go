package analysis

import (
	"context"

	"main/src/domain"
)

func basicLatency(ctx context.Context, conf domain.SimConfig, analysisTFs []analysisTF) ([]analysisTF, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for i := 0; i < len(analysisTFs); i++ {
			analysisTFs[i].Basic = calcBasicLatency(conf, analysisTFs[i])
		}
		return analysisTFs, nil
	}
}

func calcBasicLatency(conf domain.SimConfig, aTF analysisTF) int {
	noFlits := aTF.PacketSize
	processingDelay := len(aTF.Route) * conf.ProcessingDelay
	return noFlits + processingDelay
}
