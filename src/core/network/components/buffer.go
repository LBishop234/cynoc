package components

import (
	"main/log"
	"main/src/domain"
	"main/src/traffic/packet"
)

type buffer interface {
	totalCapacity() int
	vChanCapacity() int
	peakFlit(priority int) (packet.Flit, bool)
	popFlit(priority int) (packet.Flit, bool)
	addFlit(flit packet.Flit) error
}

type bufferImpl struct {
	bufferCap int
	vChanCap  int
	flits     map[int][]packet.Flit
}

func newBuffer(capacity, maxPriority int) (*bufferImpl, error) {
	if err := validBufferSize(capacity, maxPriority); err != nil {
		log.Log.Error().Err(err).Msg("invalid buffer size")
		return nil, err
	}

	vChanCap, err := bufferVChanCapacity(capacity, maxPriority)
	if err != nil {
		log.Log.Error().Err(err).Msg("invalid virtual channel size")
		return nil, err
	}

	log.Log.Trace().Msg("new buffer")
	return &bufferImpl{
		bufferCap: capacity,
		vChanCap:  vChanCap,
		flits:     make(map[int][]packet.Flit, 0),
	}, nil
}

func (b *bufferImpl) totalCapacity() int {
	return b.bufferCap
}

func (b *bufferImpl) vChanCapacity() int {
	return b.vChanCap
}

func (b *bufferImpl) peakFlit(priority int) (packet.Flit, bool) {
	if len(b.flits[priority]) == 0 {
		return nil, false
	}

	return b.flits[priority][0], true
}

func (b *bufferImpl) popFlit(priority int) (packet.Flit, bool) {
	flit, exists := b.peakFlit(priority)
	if !exists {
		return nil, false
	} else {
		if len(b.flits[priority]) > 1 {
			b.flits[priority] = b.flits[priority][1:]
		} else {
			b.flits[priority] = b.flits[priority][:0]
		}
		return flit, true
	}
}

func (b *bufferImpl) addFlit(flit packet.Flit) error {
	if len(b.flits[flit.Priority()]) < b.vChanCap {
		b.flits[flit.Priority()] = append(b.flits[flit.Priority()], flit)

		return nil
	} else {
		return domain.ErrBufferNoCapacity
	}
}

func validBufferSize(capacity, maxPriority int) error {
	if capacity < 1 {
		return domain.ErrInvalidParameter
	}

	if maxPriority < 1 {
		return domain.ErrInvalidParameter
	} else if capacity%maxPriority != 0 {
		return domain.ErrInvalidParameter
	}

	return nil
}

func bufferVChanCapacity(capacity, maxPriority int) (int, error) {
	if capacity%maxPriority != 0 {
		return 0, domain.ErrInvalidParameter
	}

	return capacity / maxPriority, nil
}
