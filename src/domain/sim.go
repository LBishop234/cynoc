package domain

type SimConfig struct {
	CycleLimit      int `yaml:"cycle_limit" json:"cycle_limit"`
	MaxPriority     int `yaml:"max_priority" json:"max_priority"`
	FlitSize        int `yaml:"flit_size" json:"flit_size"`
	BufferSize      int `yaml:"buffer_size" json:"buffer_size"`
	LinkBandwidth   int `yaml:"link_bandwidth" json:"link_bandwidth"`
	ProcessingDelay int `yaml:"processing_delay" json:"processing_delay"`
}
