package results

type resultParameter struct {
	name                string
	terminalStr         string
	csvStr              string
	terminalAllowedFlag bool
	reqAnalysisFlag     bool
	value               func(tf tfSimAnalysis) any
}

var parameters = []resultParameter{
	{
		name:                "Traffic Flow ID",
		terminalStr:         "T_i",
		csvStr:              "TF_ID",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.ID },
	},
	{
		name:                "Direct Interference Count",
		terminalStr:         "|S^D_i|",
		csvStr:              "Direct_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) any { return tf.DirectInterferenceCount },
	},
	{
		name:                "Indirect Interference Count",
		terminalStr:         "|S^I_i|",
		csvStr:              "Indirect_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) any { return tf.IndirectInterferenceCount },
	},
	{
		name:                "Number of Packets Routed",
		terminalStr:         "No. pkts",
		csvStr:              "Num_Packets_Routed",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.PacketsRouted },
	},
	{
		name:                "Number of Packets Exceeded Deadline",
		terminalStr:         "No. > D_i",
		csvStr:              "Num_Packets_Exceeded_Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.PacketsExceededDeadline },
	},
	{
		name:                "Best Latency",
		terminalStr:         "min",
		csvStr:              "Min_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return cleanInt(tf.BestLatency) },
	},
	{
		name:                "Mean Latency",
		terminalStr:         "mean",
		csvStr:              "Mean_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return cleanFloat(tf.MeanLatency) },
	},
	{
		name:                "Worst Latency",
		terminalStr:         "max",
		csvStr:              "Max_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return cleanInt(tf.WorstLatency) },
	},
	{
		name:                "Deadline",
		terminalStr:         "D_i",
		csvStr:              "Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.Deadline },
	},
	{
		name:                "Simulation Schedulable",
		terminalStr:         "Schedulable",
		csvStr:              "Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.Schedulable() },
	},
	{
		name:                "Jitter",
		terminalStr:         "J^R_i",
		csvStr:              "Jitter",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) any { return tf.Jitter },
	},
	{
		name:                "Jitter + Maximum Basic Network Latency",
		terminalStr:         "J^R_i + C_i",
		csvStr:              "Jitter_Plus_Basic",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) any { return tf.Jitter + tf.Basic },
	},
	{
		name:                "Jitter + Shi & Burns Network Latency",
		terminalStr:         "J^R_i + R_i",
		csvStr:              "Jitter_Plus_Shi_And_Burns",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) any { return tf.Jitter + tf.ShiAndBurns },
	},
	{
		name:                "Shi & Burns Schedulable",
		terminalStr:         "S&B Schedulable",
		csvStr:              "Shi_Burns_Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) any { return tf.AnalysisSchedulable() },
	},
}

// func newNewResults(simRes domain.FullResults, tfOrder []domain.TrafficFlowConfig) (domain.Results, error) {

// }
