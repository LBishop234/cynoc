package traffic

import (
	"math/rand"
	"path/filepath"

	"main/core/log"
	"main/domain"

	csvtag "github.com/artonge/go-csv-tag/v2"
)

type TrafficFlow interface {
	ID() string
	Src() string
	Dst() string
	Priority() int
	ReleasePeriod() int
	Deadline() int
	Jitter() int
	PacketSize() int
	ValidateAgainstConfig(conf domain.SimConfig) error
	ReleasePacket(cycle int) (bool, int)
}

type trafficFlowImpl struct {
	id string

	src string
	dst string

	priority      int
	releasePeriod int
	deadline      int
	jitter        int
	packetSize    int

	currentPeriod int
	currentJitter int
}

func LoadTrafficFlowConfig(fPath string) ([]domain.TrafficFlowConfig, error) {
	var trafficFlowConfigs []domain.TrafficFlowConfig
	var err error

	log.Log.Debug().Msg("reading traffic flows file")

	switch filepath.Ext(fPath) {
	case ".csv":
		err = csvtag.LoadFromPath(fPath, &trafficFlowConfigs)
		log.Log.Debug().Msg("read .csv traffic flows file")

	default:
		log.Log.Error().Err(domain.ErrInvalidFilepath).Str("ext", filepath.Ext(fPath)).Msg("invalid traffic file extension")
		return nil, domain.ErrInvalidFilepath
	}

	if err != nil {
		log.Log.Error().Err(err).Str("path", fPath).Msg("error loading traffic flows from file")
		return nil, err
	}

	return trafficFlowConfigs, nil
}

func TrafficFlows(conf domain.SimConfig, tfConfs []domain.TrafficFlowConfig) ([]TrafficFlow, error) {
	trafficFlows := make([]TrafficFlow, len(tfConfs))
	var err error

	for i := 0; i < len(trafficFlows); i++ {
		if trafficFlows[i], err = NewTrafficFlow(tfConfs[i]); err != nil {
			log.Log.Error().Err(err).Str("id", tfConfs[i].ID).Msg("error creating traffic flow")
			return nil, err
		}

		if err := trafficFlows[i].ValidateAgainstConfig(conf); err != nil {
			log.Log.Error().Err(err).Str("id", tfConfs[i].ID).Msg("error invalid traffic flow")
			return nil, err
		}
	}

	log.Log.Info().Msg("loaded traffic flows from file")
	return trafficFlows, nil
}

func NewTrafficFlow(conf domain.TrafficFlowConfig) (*trafficFlowImpl, error) {
	if conf.Src == "" {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Str("src", conf.Src).Msg("TrafficFlow source router is undefined")
		return nil, domain.ErrInvalidConfig
	}
	if conf.Dst == "" {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Str("dst", conf.Dst).Msg("TrafficFlow destination router is undefined")
		return nil, domain.ErrInvalidConfig
	}
	if conf.Priority < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("priority", conf.Priority).Msg("Invalid TrafficFlow priority")
		return nil, domain.ErrInvalidConfig
	}
	if conf.Period < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("period", conf.Period).Msg("Invalid TrafficFlow period")
		return nil, domain.ErrInvalidConfig
	}
	if conf.Deadline < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("deadline", conf.Deadline).Msg("Invalid TrafficFlow deadline")
		return nil, domain.ErrInvalidConfig
	}
	if conf.Jitter < 0 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("jitter", conf.Jitter).Msg("Invalid TrafficFlow jitter")
		return nil, domain.ErrInvalidConfig
	}
	if conf.PacketSize < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("packet_size", conf.PacketSize).Msg("Invalid TrafficFlow packet size")
		return nil, domain.ErrInvalidConfig
	}

	if conf.Deadline > conf.Period {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Int("deadline", conf.Deadline).Int("period", conf.Period).Msg("TrafficFlow deadline must be less than or equal period")
		return nil, domain.ErrInvalidConfig
	}

	if conf.Src == conf.Dst {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", conf.ID).Str("src", conf.Src).Str("dst", conf.Dst).Msg("TrafficFlow source and destination routers are the same")
		return nil, domain.ErrInvalidConfig
	}

	log.Log.Trace().Str("id", conf.ID).Msg("new traffic flow")
	return &trafficFlowImpl{
		id:            conf.ID,
		src:           conf.Src,
		dst:           conf.Dst,
		priority:      conf.Priority,
		releasePeriod: conf.Period,
		deadline:      conf.Deadline,
		jitter:        conf.Jitter,
		packetSize:    conf.PacketSize,
	}, nil
}

func (t *trafficFlowImpl) ID() string {
	return t.id
}

func (t *trafficFlowImpl) Src() string {
	return t.src
}

func (t *trafficFlowImpl) Dst() string {
	return t.dst
}

func (t *trafficFlowImpl) Priority() int {
	return t.priority
}

func (t *trafficFlowImpl) ReleasePeriod() int {
	return t.releasePeriod
}

func (t *trafficFlowImpl) Deadline() int {
	return t.deadline
}

func (t *trafficFlowImpl) Jitter() int {
	return t.jitter
}

func (t *trafficFlowImpl) PacketSize() int {
	return t.packetSize
}

func (t *trafficFlowImpl) ValidateAgainstConfig(conf domain.SimConfig) error {
	if t.priority > conf.MaxPriority {
		log.Log.Error().Str("id", t.id).
			Int("priority", t.priority).Int("max_priority", conf.MaxPriority).
			Msg("traffic flow priority exceeds max priority")
		return domain.ErrInvalidConfig
	}
	return nil
}

func (t *trafficFlowImpl) ReleasePacket(cycle int) (bool, int) {
	if cycle%t.releasePeriod == 0 {
		t.currentPeriod = cycle
		t.currentJitter = rand.Intn(t.jitter + 1)
	}

	if cycle == t.currentPeriod+t.currentJitter {
		return true, t.currentPeriod
	} else {
		return false, t.currentPeriod
	}
}
