package components

import (
	"errors"
	"sync"

	"main/domain"
	"main/log"
	"main/traffic/packet"

	"github.com/google/uuid"
)

type Router interface {
	NodeID() domain.NodeID
	RegisterInputPort(connection Connection) error
	RegisterOutputPort(connection Connection) error
	SetNetworkInterface(netIntfc NetworkInterface) error

	UpdateOutputMap()
	UpdateOutputPortsCredit() error
	ReadFromInputPorts() error
	RouteBufferedFlits() error
}

type routerImpl struct {
	// Core Attributes
	nodeID      domain.NodeID
	inputPorts  []inputPort
	outputPorts []outputPort

	outputMapSync sync.Once
	outputMap     map[domain.NodeID]outputPort

	// Configuration Constants
	routingAlg      domain.RoutingAlgorithm
	bufferSize      int
	flitSize        int
	maxPriority     int
	processingDelay int

	// Internal Operation
	headerFlitsProcessings map[uuid.UUID]int
	packetsNextRouter      map[uuid.UUID]domain.NodeID
}

type RouterConfig struct {
	domain.NodeID
	domain.SimConfig
}

func newRouter(conf RouterConfig) (*routerImpl, error) {
	if err := validBufferSize(conf.BufferSize, conf.MaxPriority); err != nil {
		return nil, err
	}

	if conf.ProcessingDelay < 1 {
		return nil, errors.Join(domain.ErrInvalidParameter, errors.New("router processing delay less then 1"))
	}

	log.Log.Trace().Str("id", conf.ID).Msg("new router")
	return &routerImpl{
		nodeID:      conf.NodeID,
		inputPorts:  make([]inputPort, 0),
		outputPorts: make([]outputPort, 0),

		outputMap: make(map[domain.NodeID]outputPort),

		routingAlg:      conf.RoutingAlgorithm,
		bufferSize:      conf.BufferSize,
		flitSize:        conf.FlitSize,
		maxPriority:     conf.MaxPriority,
		processingDelay: conf.ProcessingDelay,

		headerFlitsProcessings: make(map[uuid.UUID]int),
		packetsNextRouter:      make(map[uuid.UUID]domain.NodeID),
	}, nil
}

func (r *routerImpl) NodeID() domain.NodeID {
	return r.nodeID
}

func (r *routerImpl) RegisterInputPort(conn Connection) error {
	buff, err := newBuffer(r.bufferSize, r.maxPriority)
	if err != nil {
		return err
	}

	port, err := newInputPort(conn, buff)
	if err != nil {
		return err
	}

	conn.SetDstRouter(r.NodeID())

	r.inputPorts = append(r.inputPorts, port)

	return nil
}

func (r *routerImpl) RegisterOutputPort(conn Connection) error {
	port, err := newOutputPort(conn, r.maxPriority)
	if err != nil {
		return err
	}
	conn.SetSrcRouter(r.NodeID())

	r.outputPorts = append(r.outputPorts, port)

	return nil
}

func (r *routerImpl) UpdateOutputMap() {
	r.outputMapSync.Do(func() {
		r.outputMap = make(map[domain.NodeID]outputPort, len(r.outputPorts))
		for i := 0; i < len(r.outputPorts); i++ {
			r.outputMap[r.outputPorts[i].connection().GetDstRouter()] = r.outputPorts[i]
		}
	})
}

func (r *routerImpl) SetNetworkInterface(netIntfc NetworkInterface) error {
	if netIntfc == nil {
		return domain.ErrNilParameter
	}

	inConn, err := NewConnection(r.maxPriority)
	if err != nil {
		return err
	}

	inConn.SetSrcRouter(netIntfc.NodeID())
	inConn.SetDstRouter(r.NodeID())

	if err := netIntfc.SetOutputPort(inConn); err != nil {
		return err
	}
	if err := r.RegisterInputPort(inConn); err != nil {
		return err
	}

	outConn, err := NewConnection(r.maxPriority)
	if err != nil {
		return err
	}

	outConn.SetSrcRouter(r.NodeID())
	outConn.SetDstRouter(netIntfc.NodeID())

	if err := netIntfc.SetInputPort(outConn); err != nil {
		return err
	}
	if err := r.RegisterOutputPort(outConn); err != nil {
		return err
	}

	return nil
}

func (r *routerImpl) UpdateOutputPortsCredit() error {
	for i := 0; i < len(r.outputPorts); i++ {
		r.outputPorts[i].updateCredits()
	}

	return nil
}

func (r *routerImpl) RouteBufferedFlits() error {
	for p := 1; p <= r.maxPriority; p++ {
		for i := 0; i < len(r.inputPorts); i++ {
			if flit, exists := r.inputPorts[i].peakBuffer(p); exists {
				if err := r.arbitrateFlit(i, flit); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *routerImpl) arbitrateFlit(inputPortIndex int, flit packet.Flit) error {
	if flit.Type() == packet.HeaderFlitType {
		if headerFlit, ok := flit.(packet.HeaderFlit); ok {
			ready, err := r.processHeaderFlit(headerFlit)
			if err != nil {
				log.Log.Error().Err(err).
					Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
					Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
					Msg("error routing buffered flit")

				return err
			} else if !ready {
				return nil
			}
		} else {
			log.Log.Error().Err(domain.ErrUnknownFlitType).
				Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
				Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
				Msg("error casting header flit to packet.HeaderFlit type")

			return domain.ErrUnknownFlitType
		}
	}

	if _, exists := r.outputMap[r.packetsNextRouter[flit.PacketUUID()]]; !exists {
		if flit.Type() == packet.HeaderFlitType {
			log.Log.Error().Err(domain.ErrNoPort).
				Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
				Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
				Msg("error routing buffered flit")

			return domain.ErrNoPort
		} else {
			log.Log.Error().Err(domain.ErrMisorderedPacket).
				Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
				Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
				Msg("header flit for packet not previously processed. No output port allocated for flit")

			return domain.ErrMisorderedPacket
		}
	}

	sent, err := r.sendFlit(inputPortIndex, flit)
	if err != nil {
		log.Log.Error().Err(err).
			Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
			Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
			Msg("error sending buffered flit")

		return err
	}

	if sent {
		log.Log.Trace().
			Str("router", r.NodeID().ID).Str("packet", flit.PacketUUID().String()).
			Str("type", flit.Type().String()).Str("flit", flit.UUID().String()).
			Int("priority", flit.Priority()).
			Msg("routed buffered flit")
	}

	return nil
}

func (r *routerImpl) processHeaderFlit(flit packet.HeaderFlit) (bool, error) {
	if _, exists := r.headerFlitsProcessings[flit.UUID()]; exists {
		r.headerFlitsProcessings[flit.UUID()]++
	} else {
		r.headerFlitsProcessings[flit.UUID()] = 1
	}

	if r.headerFlitsProcessings[flit.UUID()] >= r.processingDelay {
		outPort, err := r.routeFlit(flit)
		if err != nil {
			return false, err
		}

		r.packetsNextRouter[flit.PacketUUID()] = outPort.connection().GetDstRouter()

		return true, nil
	} else {
		return false, nil
	}
}

func (r *routerImpl) sendFlit(inputPortIndex int, flit packet.Flit) (bool, error) {
	outPort, exists := r.outputMap[r.packetsNextRouter[flit.PacketUUID()]]
	if !exists {
		return false, domain.ErrInvalidParameter
	}

	if outPort.allowedToSend(flit.Priority()) {
		flit, exists := r.inputPorts[inputPortIndex].readOutOfBuffer(flit.Priority())
		if !exists {
			return false, domain.ErrInvalidParameter
		}

		if err := outPort.sendFlit(flit); err != nil {
			return false, err
		}

		return true, nil
	} else {
		return false, nil
	}
}

func (r *routerImpl) ReadFromInputPorts() error {
	for i := 0; i < len(r.inputPorts); i++ {
		err := r.inputPorts[i].readIntoBuffer()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *routerImpl) routeFlit(flit packet.HeaderFlit) (outputPort, error) {
	route := flit.Route()
	for i := 0; i < len(route); i++ {
		if route[i] == r.NodeID() {
			if i == len(route)-1 {
				return r.outputMap[route[i]], nil
			} else {
				return r.outputMap[route[i+1]], nil
			}
		}
	}

	return nil, domain.ErrNoPort
}
