package domain

import (
	"errors"
)

var (
	ErrInvalidConfig = errors.New("invalid config")

	ErrInvalidFilepath = errors.New("invalid filepath")

	ErrNilParameter     = errors.New("nil parameter passed to function")
	ErrInvalidParameter = errors.New("invalid parameter passed to function")

	ErrUnknownRoutingAlgorithm = errors.New("unknown routing algorithm")

	ErrBufferNoCapacity = errors.New("buffer has reached capacity")

	ErrPortNoCredit = errors.New("port has no credit available for flit")

	ErrUnknownFlitType = errors.New("Unknown flit type")

	ErrNoPort           = errors.New("router found no port for packet")
	ErrMisorderedPacket = errors.New("router received misordered packet flits")

	ErrInvalidTopology = errors.New("invalid topology")

	ErrMissingNetworkInterface = errors.New("missing network interface")
	ErrMissingRouter           = errors.New("missing router")

	ErrPacketsNotEqual = errors.New("packets are not equal")

	ErrFlitAlreadySet = errors.New("header or tail flit already set")
	ErrFlitUnset      = errors.New("header or tail flit not set")

	ErrMissingTrafficFlow = errors.New("missing traffic flow")
)
