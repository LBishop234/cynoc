package simulation

import (
	"math"

	"main/log"
	"main/src/traffic/packet"

	"github.com/google/uuid"
)

type records struct {
	transmittedByTF map[string]map[uuid.UUID]transmittedPacket
	arrivedByTF     map[string]map[uuid.UUID]arrivedPacket
}

type transmittedPacket struct {
	generationCycle   float64
	transmissionCycle float64
	packet            packet.Packet
}

type arrivedPacket struct {
	transmittedPacket
	receivedCycle float64
}

func newRecords() *records {
	return &records{
		transmittedByTF: make(map[string]map[uuid.UUID]transmittedPacket),
		arrivedByTF:     make(map[string]map[uuid.UUID]arrivedPacket),
	}
}

func (r *records) recordTransmittedPacket(generationCycle, transmissionCycle int, pkt packet.Packet) {
	if _, exists := r.transmittedByTF[pkt.TrafficFlowID()]; !exists {
		r.transmittedByTF[pkt.TrafficFlowID()] = make(map[uuid.UUID]transmittedPacket)
	}

	r.transmittedByTF[pkt.TrafficFlowID()][pkt.UUID()] = transmittedPacket{
		generationCycle:   float64(generationCycle),
		transmissionCycle: float64(transmissionCycle),
		packet:            pkt,
	}
	log.Log.Trace().Str("packet", pkt.UUID().String()).Msg("recording transmitted packet")
}

func (r *records) recordArrivedPacket(cycle int, pkt packet.Packet) {
	if _, exists := r.arrivedByTF[pkt.TrafficFlowID()]; !exists {
		r.arrivedByTF[pkt.TrafficFlowID()] = make(map[uuid.UUID]arrivedPacket)
	}

	if outstandingPkt, exists := r.transmittedByTF[pkt.TrafficFlowID()][pkt.UUID()]; exists {
		if err := packet.EqualPackets(outstandingPkt.packet, pkt); err != nil {
			log.Log.Error().Err(err).Str("packet", pkt.UUID().String()).Msg("packet did not match outstanding packet")
		}

		r.arrivedByTF[pkt.TrafficFlowID()][pkt.UUID()] = arrivedPacket{
			transmittedPacket: outstandingPkt,
			receivedCycle:     float64(cycle),
		}

		delete(r.transmittedByTF[pkt.TrafficFlowID()], pkt.UUID())

		log.Log.Trace().Str("packet", pkt.UUID().String()).Msg("recording arrived packet")
	} else {
		log.Log.Error().Str("packet", pkt.UUID().String()).Msg("no matching transmitted packet found")
	}
}

func (r *records) noTransmitted() int {
	count := 0
	for tfID := range r.transmittedByTF {
		count += r.noTransmittedByTF(tfID)
	}
	return count
}

func (r *records) noTransmittedByTF(tfID string) int {
	return len(r.transmittedByTF[tfID]) + len(r.arrivedByTF[tfID])
}

func (r *records) noArrived() int {
	count := 0
	for tfID := range r.arrivedByTF {
		count += r.noArrivedByTF(tfID)
	}
	return count
}

func (r *records) noArrivedByTF(tfID string) int {
	return len(r.arrivedByTF[tfID])
}

func (r *records) noLost() int {
	return r.noTransmitted() - r.noArrived()
}

func (r *records) noLostByTF(tfID string) int {
	return r.noTransmittedByTF(tfID) - r.noArrivedByTF(tfID)
}

func (r *records) noExceededDeadline() int {
	count := 0

	for tfID := range r.arrivedByTF {
		count += r.noExceededDeadlineByTF(tfID)
	}

	return count
}

func (r *records) noExceededDeadlineByTF(tfID string) int {
	count := 0

	for _, pkt := range r.arrivedByTF[tfID] {
		if !arrivedPacketInDeadline(pkt) {
			count++
		}
	}

	return count
}

func (r *records) meanLatency() float64 {
	var totalLatency float64

	for tfID := range r.arrivedByTF {
		for uuid := range r.arrivedByTF[tfID] {
			totalLatency += arrivedPacketLatency(r.arrivedByTF[tfID][uuid])
		}
	}

	return totalLatency / float64(r.noArrived())
}

func (r *records) meanLatencyByTF(tfID string) float64 {
	var totalLatency float64

	for uuid := range r.arrivedByTF[tfID] {
		totalLatency += arrivedPacketLatency(r.arrivedByTF[tfID][uuid])
	}

	return totalLatency / float64(r.noArrivedByTF(tfID))
}

func (r *records) bestLatency() int {
	var bestLatency int = math.MaxInt

	for tfID := range r.arrivedByTF {
		tfLatency := r.bestLatencyByTF(tfID)
		if tfLatency < bestLatency {
			bestLatency = tfLatency
		}
	}

	return bestLatency
}

func (r *records) bestLatencyByTF(tfID string) int {
	var bestLatency int = math.MaxInt

	for uuid := range r.arrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.arrivedByTF[tfID][uuid]))
		if latency < bestLatency {
			bestLatency = latency
		}
	}

	return bestLatency
}

func (r *records) worstLatency() int {
	var worstLatency int = math.MinInt

	for tfID := range r.arrivedByTF {
		tfLatency := r.worstLatencyByTF(tfID)
		if tfLatency > worstLatency {
			worstLatency = tfLatency
		}
	}

	return worstLatency
}

func (r *records) worstLatencyByTF(tfID string) int {
	var worstLatency int = math.MinInt

	for uuid := range r.arrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.arrivedByTF[tfID][uuid]))
		if latency > worstLatency {
			worstLatency = latency
		}
	}

	return worstLatency
}

func arrivedPacketLatency(pkt arrivedPacket) float64 {
	return math.Round(pkt.receivedCycle - pkt.generationCycle + 1)
}

func arrivedPacketInDeadline(pkt arrivedPacket) bool {
	return arrivedPacketLatency(pkt) < float64(pkt.packet.Deadline())
}
