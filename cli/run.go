package cli

import (
	"fmt"

	"main/log"
	"main/src/config"
	"main/src/core"
	"main/src/core/results"
	"main/src/topology"
	"main/src/traffic"

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
		conf = ApplyConfigOverrides(cliCtx, conf)

		top, err := topology.ReadTopology(confArgs.TopologyPath)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading topology")
		}

		trafficFlowConfigs, err := traffic.LoadTrafficFlowConfig(confArgs.TrafficPath)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error reading traffic flows file")
		}

		resultsSet, err := core.Run(conf, top, trafficFlowConfigs, analysisArgs.Analysis)
		if err != nil {
			log.Log.Fatal().Err(err).Msg("error running simulator")
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
	outputArgs := OutputArgs(cliCtx)

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
