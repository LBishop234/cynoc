package simulation

import (
	"math"

	"main/log"
	"main/src/traffic/packet"

	"github.com/google/uuid"
)

type Records struct {
	TransmittedByTF map[string]map[uuid.UUID]transmittedPacket
	ArrivedByTF     map[string]map[uuid.UUID]arrivedPacket
}

type transmittedPacket struct {
	GenerationCycle   float64
	TransmissionCycle float64
	Packet            packet.Packet
}

type arrivedPacket struct {
	transmittedPacket
	ReceivedCycle float64
}

func newRecords() *Records {
	return &Records{
		TransmittedByTF: make(map[string]map[uuid.UUID]transmittedPacket),
		ArrivedByTF:     make(map[string]map[uuid.UUID]arrivedPacket),
	}
}

func (r *Records) recordTransmittedPacket(generationCycle, transmissionCycle int, pkt packet.Packet) {
	if _, exists := r.TransmittedByTF[pkt.TrafficFlowID()]; !exists {
		r.TransmittedByTF[pkt.TrafficFlowID()] = make(map[uuid.UUID]transmittedPacket)
	}

	r.TransmittedByTF[pkt.TrafficFlowID()][pkt.UUID()] = transmittedPacket{
		GenerationCycle:   float64(generationCycle),
		TransmissionCycle: float64(transmissionCycle),
		Packet:            pkt,
	}
	log.Log.Trace().Str("packet", pkt.UUID().String()).Msg("recording transmitted packet")
}

func (r *Records) recordArrivedPacket(cycle int, pkt packet.Packet) {
	if _, exists := r.ArrivedByTF[pkt.TrafficFlowID()]; !exists {
		r.ArrivedByTF[pkt.TrafficFlowID()] = make(map[uuid.UUID]arrivedPacket)
	}

	if outstandingPkt, exists := r.TransmittedByTF[pkt.TrafficFlowID()][pkt.UUID()]; exists {
		if err := packet.EqualPackets(outstandingPkt.Packet, pkt); err != nil {
			log.Log.Error().Err(err).Str("packet", pkt.UUID().String()).Msg("packet did not match outstanding packet")
		}

		r.ArrivedByTF[pkt.TrafficFlowID()][pkt.UUID()] = arrivedPacket{
			transmittedPacket: outstandingPkt,
			ReceivedCycle:     float64(cycle),
		}

		delete(r.TransmittedByTF[pkt.TrafficFlowID()], pkt.UUID())

		log.Log.Trace().Str("packet", pkt.UUID().String()).Msg("recording arrived packet")
	} else {
		log.Log.Error().Str("packet", pkt.UUID().String()).Msg("no matching transmitted packet found")
	}
}

func (r *Records) noTransmitted() int {
	count := 0
	for tfID := range r.TransmittedByTF {
		count += r.noTransmittedByTF(tfID)
	}
	return count
}

func (r *Records) noTransmittedByTF(tfID string) int {
	return len(r.TransmittedByTF[tfID]) + len(r.ArrivedByTF[tfID])
}

func (r *Records) noArrived() int {
	count := 0
	for tfID := range r.ArrivedByTF {
		count += r.noArrivedByTF(tfID)
	}
	return count
}

func (r *Records) noArrivedByTF(tfID string) int {
	return len(r.ArrivedByTF[tfID])
}

func (r *Records) noLost() int {
	return r.noTransmitted() - r.noArrived()
}

func (r *Records) noLostByTF(tfID string) int {
	return r.noTransmittedByTF(tfID) - r.noArrivedByTF(tfID)
}

func (r *Records) noExceededDeadline() int {
	count := 0

	for tfID := range r.ArrivedByTF {
		count += r.noExceededDeadlineByTF(tfID)
	}

	return count
}

func (r *Records) noExceededDeadlineByTF(tfID string) int {
	count := 0

	for _, pkt := range r.ArrivedByTF[tfID] {
		if !arrivedPacketInDeadline(pkt) {
			count++
		}
	}

	return count
}

func (r *Records) meanLatency() float64 {
	var totalLatency float64

	for tfID := range r.ArrivedByTF {
		for uuid := range r.ArrivedByTF[tfID] {
			totalLatency += arrivedPacketLatency(r.ArrivedByTF[tfID][uuid])
		}
	}

	return totalLatency / float64(r.noArrived())
}

func (r *Records) meanLatencyByTF(tfID string) float64 {
	var totalLatency float64

	for uuid := range r.ArrivedByTF[tfID] {
		totalLatency += arrivedPacketLatency(r.ArrivedByTF[tfID][uuid])
	}

	return totalLatency / float64(r.noArrivedByTF(tfID))
}

func (r *Records) bestLatency() int {
	var bestLatency int = math.MaxInt

	for tfID := range r.ArrivedByTF {
		tfLatency := r.bestLatencyByTF(tfID)
		if tfLatency < bestLatency {
			bestLatency = tfLatency
		}
	}

	return bestLatency
}

func (r *Records) bestLatencyByTF(tfID string) int {
	var bestLatency int = math.MaxInt

	for uuid := range r.ArrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.ArrivedByTF[tfID][uuid]))
		if latency < bestLatency {
			bestLatency = latency
		}
	}

	return bestLatency
}

func (r *Records) worstLatency() int {
	var worstLatency int = math.MinInt

	for tfID := range r.ArrivedByTF {
		tfLatency := r.worstLatencyByTF(tfID)
		if tfLatency > worstLatency {
			worstLatency = tfLatency
		}
	}

	return worstLatency
}

func (r *Records) worstLatencyByTF(tfID string) int {
	var worstLatency int = math.MinInt

	for uuid := range r.ArrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.ArrivedByTF[tfID][uuid]))
		if latency > worstLatency {
			worstLatency = latency
		}
	}

	return worstLatency
}

func arrivedPacketLatency(pkt arrivedPacket) float64 {
	return math.Round(pkt.ReceivedCycle - pkt.GenerationCycle + 1)
}

func arrivedPacketInDeadline(pkt arrivedPacket) bool {
	return arrivedPacketLatency(pkt) < float64(pkt.Packet.Deadline())
}
