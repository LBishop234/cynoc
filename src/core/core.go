package core

import (
	"context"
	"sync"

	"main/log"
	"main/src/core/analysis"
	"main/src/core/network"
	"main/src/core/results"
	"main/src/core/simulation"
	"main/src/domain"
	"main/src/topology"
	"main/src/traffic"
)

func Run(conf domain.SimConfig, top *topology.Topology, trafficConf []domain.TrafficFlowConfig, runAnalysisFlag bool) (domain.Results, error) {
	network, err := network.NewNetwork(
		top,
		conf,
	)
	if err != nil {
		log.Log.Error().Err(err).Msg("error reading topology file")
		return nil, err
	}

	trafficFlows, err := traffic.TrafficFlows(conf, trafficConf)
	if err != nil {
		log.Log.Fatal().Err(err).Msg("error constructing traffic flows")
	}

	var wg sync.WaitGroup
	analysisResultsChan := make(chan analysis.AnalysisResults, 1)
	analysisErrChan := make(chan error, 1)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	if runAnalysisFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log.Log.Info().Msg("Running analysis")

			analysisResults, err := analysis.Analysis(
				ctx,
				conf,
				top,
				trafficConf,
			)
			if err != nil {
				log.Log.Error().Err(err).Msg("error running analysis")
				analysisErrChan <- err
				cancelFunc()
				return
			}

			analysisResultsChan <- analysisResults

			if !analysisResults.AnalysesSchedulable() {
				log.Log.Warn().Msg("Analysis indicates the network is not schedulable")
			}

			log.Log.Info().Msg("Finished running analysis")
		}()
	}

	simResults, err := simulation.Simulate(
		ctx,
		network,
		trafficFlows,
		conf.RoutingAlgorithm,
		conf.CycleLimit,
	)
	if err != nil {
		cancelFunc()
		log.Log.Error().Err(err).Msg("error running simulation")
		return nil, err
	}

	var resultsSet domain.Results
	if runAnalysisFlag {
		wg.Wait()
		var analysisResults analysis.AnalysisResults = <-analysisResultsChan

		resultsSet, err = results.NewResultsWithAnalysis(simResults, analysisResults, trafficConf)
	} else {
		resultsSet, err = results.NewResults(simResults, trafficConf)
	}
	if err != nil {
		log.Log.Error().Err(err).Msg("error constructing results")
		return nil, err
	}

	return resultsSet, nil
}
