package simulation

import (
	"math"

	"main/src/traffic/packet"

	"github.com/rs/zerolog"
)

type Records struct {
	TransmittedByTF map[string]map[string]transmittedPacket
	ArrivedByTF     map[string]map[string]arrivedPacket

	logger zerolog.Logger
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

func newRecords(logger zerolog.Logger) *Records {
	return &Records{
		TransmittedByTF: make(map[string]map[string]transmittedPacket),
		ArrivedByTF:     make(map[string]map[string]arrivedPacket),

		logger: logger,
	}
}

func (r *Records) recordTransmittedPacket(generationCycle, transmissionCycle int, pkt packet.Packet) {
	if _, exists := r.TransmittedByTF[pkt.TrafficFlowID()]; !exists {
		r.TransmittedByTF[pkt.TrafficFlowID()] = make(map[string]transmittedPacket)
	}

	r.TransmittedByTF[pkt.TrafficFlowID()][pkt.PacketIndex()] = transmittedPacket{
		GenerationCycle:   float64(generationCycle),
		TransmissionCycle: float64(transmissionCycle),
		Packet:            pkt,
	}
	r.logger.Trace().Str("packet", pkt.PacketIndex()).Msg("recording transmitted packet")
}

func (r *Records) recordArrivedPacket(cycle int, pkt packet.Packet) {
	if _, exists := r.ArrivedByTF[pkt.TrafficFlowID()]; !exists {
		r.ArrivedByTF[pkt.TrafficFlowID()] = make(map[string]arrivedPacket)
	}

	if outstandingPkt, exists := r.TransmittedByTF[pkt.TrafficFlowID()][pkt.PacketIndex()]; exists {
		if err := packet.EqualPackets(outstandingPkt.Packet, pkt); err != nil {
			r.logger.Error().Err(err).Str("packet", pkt.PacketIndex()).Msg("packet did not match outstanding packet")
		}

		r.ArrivedByTF[pkt.TrafficFlowID()][pkt.PacketIndex()] = arrivedPacket{
			transmittedPacket: outstandingPkt,
			ReceivedCycle:     float64(cycle),
		}

		delete(r.TransmittedByTF[pkt.TrafficFlowID()], pkt.PacketIndex())

		r.logger.Trace().Str("packet", pkt.PacketIndex()).Msg("recording arrived packet")
	} else {
		r.logger.Error().Str("packet", pkt.PacketIndex()).Msg("no matching transmitted packet found")
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
		for id := range r.ArrivedByTF[tfID] {
			totalLatency += arrivedPacketLatency(r.ArrivedByTF[tfID][id])
		}
	}

	return totalLatency / float64(r.noArrived())
}

func (r *Records) meanLatencyByTF(tfID string) float64 {
	var totalLatency float64

	for id := range r.ArrivedByTF[tfID] {
		totalLatency += arrivedPacketLatency(r.ArrivedByTF[tfID][id])
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

	for id := range r.ArrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.ArrivedByTF[tfID][id]))
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

	for id := range r.ArrivedByTF[tfID] {
		latency := int(arrivedPacketLatency(r.ArrivedByTF[tfID][id]))
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
	return arrivedPacketLatency(pkt) <= float64(pkt.Packet.Deadline())
}
