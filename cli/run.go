package cli

import (
	"context"
	"fmt"
	"main/log"
	"main/src/config"
	"main/src/core/analysis"
	"main/src/core/network"
	"main/src/core/results"
	"main/src/core/simulation"
	"main/src/domain"
	"main/src/topology"
	"main/src/traffic"
	"sync"

	"github.com/urfave/cli/v2"
)

const (
	appName = "CyNoC"
)

func NewApp() *cli.App {
	app := &cli.App{}

	logArgs := LogArgs(app)
	analysisArgs := AnalysisArgs(app)
	confArgs := ConfigFilesArgs(app)
	ConfigOverridesArgs(app)
	SetupOutputArgs(app)

	app.Name = appName

	app.Action = func(cliCtx *cli.Context) error {
		initLogger(logArgs)

		conf, err := config.ReadConfig(confArgs.ConfigPath)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading config file")
		}

		conf = config.ApplyConfigOverrides(cliCtx, conf)

		top, err := topology.ReadTopology(confArgs.TopologyPath)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading topology")
		}

		network, err := network.NewNetwork(
			top,
			domain.SimConfig{
				RoutingAlgorithm: conf.RoutingAlgorithm,
				BufferSize:       conf.BufferSize,
				FlitSize:         conf.FlitSize,
				MaxPriority:      conf.MaxPriority,
				ProcessingDelay:  conf.ProcessingDelay,
			},
		)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading topology file")
		}

		trafficFlowConfigs, err := traffic.LoadTrafficFlowConfig(confArgs.TrafficPath)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading traffic flows file")
		}

		trafficFlows, err := traffic.TrafficFlows(conf, trafficFlowConfigs)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error constructing traffic flows")
		}

		var wg sync.WaitGroup
		analysisResultsChan := make(chan analysis.AnalysisResults, 1)
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()

		if analysisArgs.Analysis {
			wg.Add(1)
			go func() {
				defer wg.Done()

				log.Log.Info().Msg("Running analysis")

				analysisResults, err := analysis.Analysis(
					ctx,
					conf,
					top,
					trafficFlowConfigs,
				)
				if err != nil {
					log.Log.Error().Err(err).Msg("error running analysis")
					cancelFunc()
					return
				}

				analysisResultsChan <- analysisResults

				if !analysisResults.AnalysesSchedulable() {
					log.Log.Warn().Msg("Analysis indicates the network is not schedulable")
				}
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
			log.Log.Fatal().Err(err).Msg("error running simulation")
		}

		var resultsSet results.Results
		if analysisArgs.Analysis {
			wg.Wait()
			var analysisResults analysis.AnalysisResults = <-analysisResultsChan

			resultsSet, err = results.NewResultsWithAnalysis(simResults, analysisResults, trafficFlowConfigs)
		} else {
			resultsSet, err = results.NewResults(simResults, trafficFlowConfigs)
		}
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error constructing results")
		}

		if err := output(cliCtx, resultsSet); err != nil {
			log.Log.Fatal().Err(err).Msg("error outputting results")
		}

		return nil
	}

	return app
}

func initLogger(logConf *LogConfig) {
	var logLevel log.LogLevel
	if logConf.TraceOutput {
		logLevel = log.TRACE
	} else if logConf.DebugOutput {
		logLevel = log.DEBUG
	} else if logConf.Log {
		logLevel = log.INFO
	} else {
		logLevel = log.ERROR
	}
	log.InitLogger(logLevel)
}

func output(cliCtx *cli.Context, results results.Results) error {
	outputArgs := config.OutputArgs(cliCtx)

	if outputArgs.OutputFileFlag {
		if err := results.OutputCSV(outputArgs.OutputFilepath); err != nil {
			log.Log.Error().Err(err).Msgf("error writing traffic flow results to %s", outputArgs.OutputFilepath)
			return err
		}
	}

	if !outputArgs.NoConsoleOutput {
		str, err := results.Prettify()
		if err != nil {
			log.Log.Error().Err(err).Msg("error prettifying results")
			return err
		}

		fmt.Println(str)
	}

	return nil
}
