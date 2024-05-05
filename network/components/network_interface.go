package components

import (
	"main/domain"
	"main/log"
	"main/traffic/packet"

	"github.com/google/uuid"
)

type NetworkInterface interface {
	NodeID() domain.NodeID

	SetInputPort(conn Connection) error
	SetOutputPort(conn Connection) error

	RoutePacket(packet packet.Packet) error
	PopArrivedPackets() []packet.Packet

	TransmitPendingPackets() error
	HandleArrivingFlits() error
}

type networkInterfaceImpl struct {
	nodeID      domain.NodeID
	bufferSize  int
	flitSize    int
	maxPriority int

	pendingPackets map[int][]packet.Packet
	flitsInTransit map[int][]packet.Flit
	outputPort     outputPort

	inputPort      inputPort
	flitsArriving  map[uuid.UUID]packet.Reconstructor
	arrivedPackets []packet.Packet
}

func newNetworkInterface(nodeID domain.NodeID, bufferSize, flitSize, maxPriority int) (*networkInterfaceImpl, error) {
	if err := validBufferSize(bufferSize, maxPriority); err != nil {
		log.Log.Error().Err(err).Msg("invalid buffer size")
		return nil, err
	}

	log.Log.Trace().Str("id", nodeID.ID).Msg("new network interface")
	return &networkInterfaceImpl{
		nodeID:         nodeID,
		bufferSize:     bufferSize,
		flitSize:       flitSize,
		maxPriority:    maxPriority,
		pendingPackets: make(map[int][]packet.Packet, 0),
		flitsInTransit: make(map[int][]packet.Flit, 0),
		flitsArriving:  make(map[uuid.UUID]packet.Reconstructor),
		arrivedPackets: make([]packet.Packet, 0),
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

	buff, err := newBuffer(n.bufferSize, n.maxPriority)
	if err != nil {
		return err
	}

	n.inputPort, err = newInputPort(conn, buff)
	return err
}

func (n *networkInterfaceImpl) SetOutputPort(conn Connection) error {
	if conn == nil {
		return domain.ErrNilParameter
	}

	conn.SetSrcRouter(n.NodeID())

	var err error
	n.outputPort, err = newOutputPort(conn, n.maxPriority)
	return err
}

func (n *networkInterfaceImpl) RoutePacket(pkt packet.Packet) error {
	if pkt == nil {
		return domain.ErrNilParameter
	}

	n.pendingPackets[pkt.Priority()] = append(n.pendingPackets[pkt.Priority()], pkt)

	log.Log.Trace().Str("network_interface", n.NodeID().ID).Str("packet", pkt.UUID().String()).
		Msg("network interface received packet")

	return nil
}

func (n *networkInterfaceImpl) PopArrivedPackets() []packet.Packet {
	pkts := n.arrivedPackets
	n.arrivedPackets = n.arrivedPackets[:0]
	return pkts
}

func (n *networkInterfaceImpl) HandleArrivingFlits() error {
	if err := n.inputPort.readIntoBuffer(); err != nil {
		return err
	}

	actionFlag := true
	for actionFlag {
		actionFlag = false

		for p := 1; p <= n.maxPriority; p++ {
			flit, exists := n.inputPort.readOutOfBuffer(p)
			if exists {
				actionFlag = true
			} else {
				continue
			}

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
				log.Log.Error().Err(err).
					Str("network_interface", n.NodeID().ID).Str("packet", flit.PacketUUID().String()).
					Str("flit", flit.UUID().String()).Str("type", flit.Type().String()).
					Msg("error handling arrived flit")
				return err
			}

			log.Log.Trace().
				Str("network_interface", n.NodeID().ID).Str("packet", flit.PacketUUID().String()).
				Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
				Int("priority", flit.Priority()).
				Msg("flit arrived at network interface")
		}
	}

	return nil
}

func (n *networkInterfaceImpl) arrivedHeaderFlit(flit packet.HeaderFlit) error {
	n.flitsArriving[flit.PacketUUID()] = packet.NewReconstructor()

	if err := n.flitsArriving[flit.PacketUUID()].SetHeader(flit); err != nil {
		return err
	}

	return nil
}

func (n *networkInterfaceImpl) arrivedBodyFlit(flit packet.BodyFlit) error {
	reconstructor, exists := n.flitsArriving[flit.PacketUUID()]
	if !exists {
		return domain.ErrMisorderedPacket
	}

	if err := reconstructor.AddBody(flit); err != nil {
		return err
	}

	return nil
}

func (n *networkInterfaceImpl) arrivedTailFlit(flit packet.TailFlit) error {
	reconstructor, exists := n.flitsArriving[flit.PacketUUID()]
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
	delete(n.flitsArriving, flit.PacketUUID())

	return nil
}

func (n *networkInterfaceImpl) TransmitPendingPackets() error {
	n.outputPort.updateCredits()

	for p := 0; p <= n.maxPriority; p++ {
		n.updateFlitsInTransit(p)

		if len(n.flitsInTransit[p]) > 0 {
			for i := 0; i < len(n.flitsInTransit[p]); i++ {
				if n.outputPort.allowedToSend(n.flitsInTransit[p][0].Priority()) {
					if err := n.outputPort.sendFlit(n.flitsInTransit[p][0]); err != nil {
						log.Log.Error().Err(err).
							Str("network_interface", n.NodeID().ID).Str("packet", n.flitsInTransit[p][0].PacketUUID().String()).
							Str("flit", n.flitsInTransit[p][0].UUID().String()).Str("type", n.flitsInTransit[p][0].Type().String()).
							Msg("error sending flit")

						return err
					}

					log.Log.Trace().
						Str("network_interface", n.NodeID().ID).Str("packet", n.flitsInTransit[p][0].PacketUUID().String()).
						Str("flit", n.flitsInTransit[p][0].UUID().String()).Str("type", n.flitsInTransit[p][0].Type().String()).
						Int("priority", n.flitsInTransit[p][0].Priority()).
						Msg("flit sent from network interface")

					n.flitsInTransit[p] = n.flitsInTransit[p][1:]
					return nil
				}
			}
		}
	}

	return nil
}

func (n *networkInterfaceImpl) updateFlitsInTransit(priority int) {
	if len(n.flitsInTransit[priority]) == 0 && len(n.pendingPackets[priority]) > 0 {
		n.flitsInTransit[priority] = n.pendingPackets[priority][0].Flits(n.flitSize)
		n.pendingPackets[priority] = n.pendingPackets[priority][1:]
	}
}
