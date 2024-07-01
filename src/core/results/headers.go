package results

import "strconv"

type resultParameter struct {
	name                string
	terminalStr         string
	csvStr              string
	terminalAllowedFlag bool
	reqAnalysisFlag     bool
	value               func(tf tfSimAnalysis) string
}

var parameters = []resultParameter{
	{
		name:                "Traffic Flow ID",
		terminalStr:         "T_i",
		csvStr:              "TF_ID",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return tf.ID },
	},
	{
		name:                "Direct Interference Count",
		terminalStr:         "|S^D_i|",
		csvStr:              "Direct_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.DirectInterferenceCount) },
	},
	{
		name:                "Indirect Interference Count",
		terminalStr:         "|S^I_i|",
		csvStr:              "Indirect_Interference_Count",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.IndirectInterferenceCount) },
	},
	{
		name:                "Number of Packets Routed",
		terminalStr:         "No. pkts",
		csvStr:              "Num_Packets_Routed",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.PacketsRouted) },
	},
	{
		name:                "Number of Packets Exceeded Deadline",
		terminalStr:         "No. > D_i",
		csvStr:              "Num_Packets_Exceeded_Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.PacketsExceededDeadline) },
	},
	{
		name:                "Best Latency",
		terminalStr:         "min",
		csvStr:              "Min_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanInt(tf.BestLatency) },
	},
	{
		name:                "Mean Latency",
		terminalStr:         "mean",
		csvStr:              "Mean_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanFloat(tf.MeanLatency) },
	},
	{
		name:                "Worst Latency",
		terminalStr:         "max",
		csvStr:              "Max_Latency",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return cleanInt(tf.WorstLatency) },
	},
	{
		name:                "Deadline",
		terminalStr:         "D_i",
		csvStr:              "Deadline",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Deadline) },
	},
	{
		name:                "Simulation Schedulable",
		terminalStr:         "Schedulable",
		csvStr:              "Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.FormatBool(tf.Schedulable()) },
	},
	{
		name:                "Jitter",
		terminalStr:         "J^R_i",
		csvStr:              "Jitter",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     false,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter) },
	},
	{
		name:                "Jitter + Maximum Basic Network Latency",
		terminalStr:         "J^R_i + C_i",
		csvStr:              "Jitter_Plus_Basic",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter + tf.Basic) },
	},
	{
		name:                "Jitter + Shi & Burns Network Latency",
		terminalStr:         "J^R_i + R_i",
		csvStr:              "Jitter_Plus_Shi_Burns",
		terminalAllowedFlag: true,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.Itoa(tf.Jitter + tf.ShiAndBurns) },
	},
	{
		name:                "Shi & Burns Schedulable",
		terminalStr:         "S&B Schedulable",
		csvStr:              "Shi_Burns_Schedulable",
		terminalAllowedFlag: false,
		reqAnalysisFlag:     true,
		value:               func(tf tfSimAnalysis) string { return strconv.FormatBool(tf.AnalysisSchedulable()) },
	},
}
