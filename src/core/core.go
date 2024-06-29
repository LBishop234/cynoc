package core

import (
	"context"
	"sync"

	"main/src/core/analysis"
	"main/src/core/network"
	"main/src/core/results"
	"main/src/core/simulation"
	"main/src/domain"
	"main/src/topology"
	"main/src/traffic"

	"github.com/rs/zerolog"
)

func Run(conf domain.SimConfig, top *topology.Topology, trafficConf []domain.TrafficFlowConfig, runAnalysisFlag bool, logger zerolog.Logger) (results.Results, error) {
	network, err := network.NewNetwork(
		top,
		conf,
		logger,
	)
	if err != nil {
		logger.Error().Err(err).Msg("error reading topology file")
		return nil, err
	}

	trafficFlows, err := traffic.TrafficFlows(conf, trafficConf)
	if err != nil {
		logger.Fatal().Err(err).Msg("error constructing traffic flows")
	}

	var wg sync.WaitGroup
	analysisResultsChan := make(chan domain.AnalysisResults, 1)
	analysisErrChan := make(chan error, 1)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	if runAnalysisFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()

			logger.Info().Msg("Running analysis")

			analysisResults, err := analysis.Analysis(
				ctx,
				conf,
				top,
				trafficConf,
			)
			if err != nil {
				logger.Error().Err(err).Msg("error running analysis")
				analysisErrChan <- err
				cancelFunc()
				return
			}

			analysisResultsChan <- analysisResults

			if !analysisResults.AnalysesSchedulable() {
				logger.Warn().Msg("Analysis indicates the network is not schedulable")
			}

			logger.Info().Msg("Finished running analysis")
		}()
	}

	simResults, err := simulation.Simulate(
		ctx,
		network,
		trafficFlows,
		conf.RoutingAlgorithm,
		conf.CycleLimit,
		logger,
	)
	if err != nil {
		cancelFunc()
		logger.Error().Err(err).Msg("error running simulation")
		return nil, err
	}

	var resultsSet results.Results
	if runAnalysisFlag {
		wg.Wait()
		var analysisResults domain.AnalysisResults = <-analysisResultsChan

		resultsSet, err = results.NewResultsWithAnalysis(simResults, analysisResults, trafficConf)
	} else {
		resultsSet, err = results.NewResults(simResults, trafficConf)
	}
	if err != nil {
		logger.Error().Err(err).Msg("error constructing results")
		return nil, err
	}

	return resultsSet, nil
}
