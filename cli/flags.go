package cli

import (
	"fmt"

	"main/src/domain"

	"github.com/urfave/cli/v2"
)

type (
	LogConfig struct {
		Log         bool
		DebugOutput bool
		TraceOutput bool
	}

	Analysis struct {
		Analysis bool
	}

	ConfigFiles struct {
		ConfigPath   string
		TopologyPath string
		TrafficPath  string
	}

	outputArgs struct {
		NoConsoleOutput bool
		OutputFileFlag  bool
		OutputFilepath  string
	}
)

func LogArgs(app *cli.App) *LogConfig {
	const category = "Logging"

	lConf := &LogConfig{}

	app.Flags = append(
		app.Flags,
		&cli.BoolFlag{
			Name:        "log",
			Usage:       "print logs to stdout",
			Value:       false,
			Destination: &lConf.Log,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "print debug output",
			Value:       false,
			Destination: &lConf.DebugOutput,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "trace",
			Usage:       "print trace output",
			Value:       false,
			Destination: &lConf.TraceOutput,
			Category:    category,
		},
	)

	return lConf
}

func AnalysisArgs(app *cli.App) *Analysis {
	const category = "Analysis"

	analysis := &Analysis{}

	app.Flags = append(
		app.Flags,
		&cli.BoolFlag{
			Name:        "analysis",
			Aliases:     []string{"a"},
			Usage:       "enable analysis",
			Value:       false,
			Destination: &analysis.Analysis,
			Category:    category,
		},
	)

	return analysis
}

func ConfigFilesArgs(app *cli.App) *ConfigFiles {
	const category = "Configuration Files"

	cConf := &ConfigFiles{}

	app.Flags = append(
		app.Flags,
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "load configuration from `FILE`",
			Destination: &cConf.ConfigPath,
			Required:    true,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "topology",
			Aliases:     []string{"t"},
			Usage:       "load topology from `FILE`",
			Destination: &cConf.TopologyPath,
			Required:    true,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "traffic",
			Aliases:     []string{"tr"},
			Usage:       "load traffic from `FILE`",
			Destination: &cConf.TrafficPath,
			Required:    true,
			Category:    category,
		},
	)

	return cConf
}

const (
	overrideCyclesFlag     = "cycle_limit"
	overideMaxPriorityFlag = "max_priority"
	overrideBufferSizeFlag = "buffer_size"
	processingDelayFlag    = "processing_delay"
)

func ConfigOverridesArgs(app *cli.App) {
	const category = "Configuration Overrides"
	const usageBaseStr = "override %s's configuration file value with `VALUE`"

	app.Flags = append(
		app.Flags,
		&cli.IntFlag{
			Name:        overrideCyclesFlag,
			Aliases:     []string{"cy"},
			Usage:       fmt.Sprintf(usageBaseStr, overrideCyclesFlag),
			Category:    category,
			DefaultText: "no-op when unset",
		},
		&cli.StringFlag{
			Name:        overideMaxPriorityFlag,
			Aliases:     []string{"mp"},
			Usage:       fmt.Sprintf(usageBaseStr, overideMaxPriorityFlag),
			Category:    category,
			DefaultText: "no-op when unset",
		},
		&cli.IntFlag{
			Name:        overrideBufferSizeFlag,
			Aliases:     []string{"bs"},
			Usage:       fmt.Sprintf(usageBaseStr, overrideBufferSizeFlag),
			Category:    category,
			DefaultText: "no-op when unset",
		},
		&cli.IntFlag{
			Name:        processingDelayFlag,
			Aliases:     []string{"pd"},
			Usage:       fmt.Sprintf(usageBaseStr, processingDelayFlag),
			Category:    category,
			DefaultText: "no-op when unset",
		},
	)
}

func ApplyConfigOverrides(ctx *cli.Context, conf domain.SimConfig) domain.SimConfig {
	if ctx.IsSet(overrideCyclesFlag) {
		conf.CycleLimit = ctx.Int(overrideCyclesFlag)
	}
	if ctx.IsSet(overideMaxPriorityFlag) {
		conf.MaxPriority = ctx.Int(overideMaxPriorityFlag)
	}
	if ctx.IsSet(overrideBufferSizeFlag) {
		conf.BufferSize = ctx.Int(overrideBufferSizeFlag)
	}
	if ctx.IsSet(processingDelayFlag) {
		conf.ProcessingDelay = ctx.Int(processingDelayFlag)
	}
	return conf
}

var outputFileFlag = "results-csv"

func SetupOutputArgs(app *cli.App) {
	const category = "Output"

	app.Flags = append(
		app.Flags,
		&cli.BoolFlag{
			Name:     "no-console-output",
			Aliases:  []string{"nco"},
			Usage:    "disable console output",
			Category: category,
		},
		&cli.StringFlag{
			Name:     outputFileFlag,
			Aliases:  []string{"csv"},
			Usage:    "store output results csv to `FILE`",
			Category: category,
		},
	)
}

func OutputArgs(ctx *cli.Context) outputArgs {
	var oArgs outputArgs

	oArgs.NoConsoleOutput = ctx.Bool("no-console-output")

	if ctx.IsSet(outputFileFlag) {
		oArgs.OutputFileFlag = true
		oArgs.OutputFilepath = ctx.String(outputFileFlag)
	}

	return oArgs
}
