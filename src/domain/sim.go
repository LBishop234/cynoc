package domain

type SimConfig struct {
	CycleLimit      int `yaml:"cycle_limit" json:"cycle_limit"`
	MaxPriority     int `yaml:"max_priority" json:"max_priority"`
	BufferSize      int `yaml:"buffer_size" json:"buffer_size"`
	ProcessingDelay int `yaml:"processing_delay" json:"processing_delay"`
}
