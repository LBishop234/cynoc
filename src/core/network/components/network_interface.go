package components

import (
	"main/src/domain"
	"main/src/traffic/packet"

	"github.com/rs/zerolog"
)

type NetworkInterface interface {
	NodeID() domain.NodeID

	SetInputPort(conn Connection) error
	SetOutputPort(conn Connection) error

	RoutePacket(cycle int, packet packet.Packet) error
	PopArrivedPackets(cycle int) []packet.Packet

	TransmitPendingPackets(cycle int) error
	HandleArrivingFlits(cycle int) error
}

type networkInterfaceImpl struct {
	// Core Attributes
	nodeID      domain.NodeID
	bufferSize  int
	flitSize    int
	maxPriority int

	flitsInTransit map[int][]packet.Flit
	outputPort     outputPort

	inputPort      inputPort
	flitsArriving  map[string]packet.Reconstructor
	arrivedPackets []packet.Packet

	// Utility
	logger zerolog.Logger
}

func newNetworkInterface(nodeID domain.NodeID, bufferSize, flitSize, maxPriority int, logger zerolog.Logger) (*networkInterfaceImpl, error) {
	if err := validBufferSize(bufferSize, maxPriority); err != nil {
		logger.Error().Err(err).Msg("invalid buffer size")
		return nil, err
	}

	logger.Trace().Str("id", nodeID.ID).Msg("new network interface")
	return &networkInterfaceImpl{
		nodeID:         nodeID,
		bufferSize:     bufferSize,
		flitSize:       flitSize,
		maxPriority:    maxPriority,
		flitsInTransit: make(map[int][]packet.Flit),
		flitsArriving:  make(map[string]packet.Reconstructor),
		arrivedPackets: make([]packet.Packet, 0),

		logger: logger.With().Str("component", "network_interface").Str("node_id", nodeID.ID).Logger(),
	}, nil
}

func (n *networkInterfaceImpl) NodeID() domain.NodeID {
	return n.nodeID
}

func (n *networkInterfaceImpl) SetInputPort(conn Connection) error {
	if conn == nil {
		return domain.ErrNilParameter
	}

	conn.SetDstRouter(n.NodeID())

	buff, err := newBuffer(n.bufferSize, n.maxPriority, n.logger)
	if err != nil {
		return err
	}

	n.inputPort, err = newInputPort(conn, buff, n.logger)
	return err
}

func (n *networkInterfaceImpl) SetOutputPort(conn Connection) error {
	if conn == nil {
		return domain.ErrNilParameter
	}

	conn.SetSrcRouter(n.NodeID())

	var err error
	n.outputPort, err = newOutputPort(conn, n.maxPriority, n.logger)
	return err
}

func (n *networkInterfaceImpl) RoutePacket(cycle int, pkt packet.Packet) error {
	logger := n.logger.With().Int("cycle", cycle).Logger()

	if pkt == nil {
		return domain.ErrNilParameter
	}

	logger.Trace().Str("packet", pkt.ID()).Msg("network interface received packet")

	flits := pkt.Flits(n.flitSize)
	for i := 0; i < len(flits); i++ {
		logger.Trace().
			Str("flit", flits[i].ID()).Str("type", flits[i].Type().String()).
			Msg("flit created at network interface")
		n.flitsInTransit[pkt.Priority()] = append(n.flitsInTransit[pkt.Priority()], flits[i])
	}

	return nil
}

func (n *networkInterfaceImpl) PopArrivedPackets(cycle int) []packet.Packet {
	pkts := n.arrivedPackets
	n.arrivedPackets = n.arrivedPackets[:0]
	return pkts
}

func (n *networkInterfaceImpl) HandleArrivingFlits(cycle int) error {
	logger := n.logger.With().Int("cycle", cycle).Logger()

	if err := n.inputPort.readIntoBuffer(cycle); err != nil {
		return err
	}

	actionFlag := true
	for actionFlag {
		actionFlag = false

		for p := 1; p <= n.maxPriority; p++ {
			for b := 0; b < n.bufferSize; b++ {
				flit, exists := n.inputPort.readOutOfBuffer(cycle, p)

				if !exists {
					break
				}

				actionFlag = true

				var err error
				if headerFlit, ok := flit.(packet.HeaderFlit); ok {
					err = n.arrivedHeaderFlit(headerFlit)
				} else if bodyFlit, ok := flit.(packet.BodyFlit); ok {
					err = n.arrivedBodyFlit(bodyFlit)
				} else if tailFlit, ok := flit.(packet.TailFlit); ok {
					err = n.arrivedTailFlit(tailFlit)
				} else {
					return domain.ErrUnknownFlitType
				}
				if err != nil {
					logger.Error().Err(err).
						Str("flit", flit.ID()).Str("type", flit.Type().String()).
						Msg("error handling arrived flit")
					return err
				}

				logger.Trace().
					Str("flit", flit.ID()).Str("type", flit.Type().String()).
					Msg("flit arrived at network interface")
			}
		}
	}

	return nil
}

func (n *networkInterfaceImpl) arrivedHeaderFlit(flit packet.HeaderFlit) error {
	_, exists := n.flitsArriving[flit.PacketID()]
	if exists {
		return domain.ErrMisorderedPacket
	}

	reconstructor, err := packet.NewReconstructor(flit, n.logger)
	if err != nil {
		n.logger.Error().Err(err).Str("packet", flit.PacketID()).Msg("error creating reconstructor")
		return err
	}

	n.flitsArriving[flit.PacketID()] = reconstructor
	return nil
}

func (n *networkInterfaceImpl) arrivedBodyFlit(flit packet.BodyFlit) error {
	reconstructor, exists := n.flitsArriving[flit.PacketID()]
	if !exists {
		return domain.ErrMisorderedPacket
	}

	if err := reconstructor.AddBody(flit); err != nil {
		return err
	}

	return nil
}

func (n *networkInterfaceImpl) arrivedTailFlit(flit packet.TailFlit) error {
	reconstructor, exists := n.flitsArriving[flit.PacketID()]
	if !exists {
		return domain.ErrMisorderedPacket
	}

	if err := reconstructor.SetTail(flit); err != nil {
		return err
	}

	packet, err := reconstructor.Reconstruct()
	if err != nil {
		return err
	}

	n.arrivedPackets = append(n.arrivedPackets, packet)
	delete(n.flitsArriving, flit.PacketID())

	return nil
}

func (n *networkInterfaceImpl) TransmitPendingPackets(cycle int) error {
	logger := n.logger.With().Int("cycle", cycle).Logger()

	n.outputPort.updateCredits()

	for p := 1; p <= n.maxPriority; p++ {
		for len(n.flitsInTransit[p]) > 0 && n.outputPort.allowedToSend(n.flitsInTransit[p][0].Priority()) {
			if err := n.outputPort.sendFlit(cycle, n.flitsInTransit[p][0]); err != nil {
				logger.Error().Err(err).
					Str("flit", n.flitsInTransit[p][0].ID()).Str("type", n.flitsInTransit[p][0].Type().String()).
					Msg("error sending flit")

				return err
			}

			logger.Trace().
				Str("flit", n.flitsInTransit[p][0].ID()).Str("type", n.flitsInTransit[p][0].Type().String()).
				Msg("flit sent from network interface")

			n.flitsInTransit[p] = n.flitsInTransit[p][1:]
		}
	}

	return nil
}
