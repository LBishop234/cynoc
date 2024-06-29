package domain

import (
	"errors"
	"strings"
)

type TrafficFlowConfig struct {
	ID         string `csv:"id"`
	Priority   int    `csv:"priority"`
	Period     int    `csv:"period"`
	Deadline   int    `csv:"deadline"`
	Jitter     int    `csv:"jitter"`
	PacketSize int    `csv:"packet_size"`
	Route      string `csv:"route"`
}

func (t *TrafficFlowConfig) RouteArray() ([]string, error) {
	if len(t.Route) < 2 {
		return nil, errors.Join(ErrInvalidConfig, ErrInvalidRoute)
	}

	str := strings.Replace(t.Route, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	return strings.Split(str, ","), nil
}
