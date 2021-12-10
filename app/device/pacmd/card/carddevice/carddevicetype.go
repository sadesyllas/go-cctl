package carddevice

import (
	"fmt"
	"strings"
)

type CardDeviceType uint64

const (
	Source CardDeviceType = iota + 1
	Sink
)

type ParseableCardDeviceType string

func (value ParseableCardDeviceType) Parse() (cardDeviceType CardDeviceType, err error) {
	switch strings.ToLower(string(value)) {
	case "source":
		cardDeviceType = Source
	case "sink":
		cardDeviceType = Sink
	default:
		err = fmt.Errorf("Invalid card device type: %v\n", value)
	}

	return
}

func (value CardDeviceType) String() string {
	switch value {
	case Source:
		return "source"
	case Sink:
		return "sink"
	}

	panic("unreachable")
}
