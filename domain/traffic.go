package domain

type TrafficFlowConfig struct {
	ID         string `csv:"id"`
	Src        string `csv:"src"`
	Dst        string `csv:"dst"`
	Priority   int    `csv:"priority"`
	Period     int    `csv:"period"`
	Deadline   int    `csv:"deadline"`
	Jitter     int    `csv:"jitter"`
	PacketSize int    `csv:"packet_size"`
}
