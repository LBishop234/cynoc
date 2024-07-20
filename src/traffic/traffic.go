package traffic

import (
	"encoding/hex"
	"math/rand"
	"path/filepath"
	"strconv"

	"main/log"
	"main/src/domain"
	"main/src/traffic/packet"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/rs/zerolog"
)

type TrafficFlow interface {
	ID() string
	Priority() int
	ReleasePeriod() int
	Deadline() int
	Jitter() int
	PacketSize() int
	Route() []string
	ReleasePacket(cycle int, trafficFlow TrafficFlow, route domain.Route, logger zerolog.Logger) (bool, packet.Packet, int)
}

type trafficFlowImpl struct {
	id string

	priority      int
	releasePeriod int
	deadline      int
	jitter        int
	packetSize    int
	route         []string

	currentPeriod int
	currentJitter int

	packetCount int
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
		if trafficFlows[i], err = NewTrafficFlow(tfConfs[i], conf); err != nil {
			log.Log.Error().Err(err).Str("id", tfConfs[i].ID).Msg("error creating traffic flow")
			return nil, err
		}
	}

	log.Log.Info().Msg("loaded traffic flows from file")
	return trafficFlows, nil
}

func NewTrafficFlow(tfConf domain.TrafficFlowConfig, conf domain.SimConfig) (*trafficFlowImpl, error) {
	if tfConf.Priority < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("priority", tfConf.Priority).Msg("Invalid TrafficFlow priority")
		return nil, domain.ErrInvalidConfig
	}
	if tfConf.Period < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("period", tfConf.Period).Msg("Invalid TrafficFlow period")
		return nil, domain.ErrInvalidConfig
	}
	if tfConf.Deadline < 1 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("deadline", tfConf.Deadline).Msg("Invalid TrafficFlow deadline")
		return nil, domain.ErrInvalidConfig
	}
	if tfConf.Jitter < 0 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("jitter", tfConf.Jitter).Msg("Invalid TrafficFlow jitter")
		return nil, domain.ErrInvalidConfig
	}
	if tfConf.PacketSize < 2 {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("packet_size", tfConf.PacketSize).Msg("Invalid TrafficFlow packet size, must be at least 2 to allow for header and tail flits")
		return nil, domain.ErrInvalidConfig
	}

	if tfConf.Deadline > (tfConf.Period - tfConf.Jitter) {
		log.Log.Error().Err(domain.ErrInvalidConfig).Str("id", tfConf.ID).Int("deadline", tfConf.Deadline).Int("period", tfConf.Period).Int("jitter", tfConf.Jitter).Msg("TrafficFlow deadline must be less than or equal period - jitter")
		return nil, domain.ErrInvalidConfig
	}

	if tfConf.Priority > conf.MaxPriority {
		log.Log.Error().Str("id", tfConf.ID).Int("priority", tfConf.Priority).Int("max_priority", conf.MaxPriority).Msg("traffic flow priority exceeds max priority")
		return nil, domain.ErrInvalidConfig
	}

	route, err := tfConf.RouteArray()
	if err != nil {
		log.Log.Error().Err(err).Str("id", tfConf.ID).Str("route", tfConf.Route).Msg("Invalid TrafficFlow route")
		return nil, err
	}

	log.Log.Trace().Str("id", tfConf.ID).Msg("new traffic flow")
	return &trafficFlowImpl{
		id:            tfConf.ID,
		priority:      tfConf.Priority,
		releasePeriod: tfConf.Period,
		deadline:      tfConf.Deadline,
		jitter:        tfConf.Jitter,
		packetSize:    tfConf.PacketSize,
		route:         route,
	}, nil
}

func (t *trafficFlowImpl) ID() string {
	return t.id
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

func (t *trafficFlowImpl) Route() []string {
	return t.route
}

func (t *trafficFlowImpl) ReleasePacket(cycle int, trafficFlow TrafficFlow, route domain.Route, logger zerolog.Logger) (bool, packet.Packet, int) {
	if cycle%t.releasePeriod == 0 {
		t.currentPeriod = cycle
		t.currentJitter = rand.Intn(t.jitter + 1)
	}

	if cycle == t.currentPeriod+t.currentJitter {
		pkt := packet.NewPacket(
			trafficFlow.ID(),
			hex.EncodeToString([]byte(strconv.Itoa(t.packetCount))),
			trafficFlow.Priority(),
			trafficFlow.Deadline(),
			route,
			trafficFlow.PacketSize(),
			logger,
		)

		t.packetCount++

		return true, pkt, t.currentPeriod
	} else {
		return false, nil, t.currentPeriod
	}
}
