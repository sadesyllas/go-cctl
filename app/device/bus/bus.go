package bus

import (
	"context"
	"fmt"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

type Bus uint64

const (
	PCI Bus = iota + 1
	Bluetooth
	USB
)

type ParseableBus string

func (value ParseableBus) Parse(ctx context.Context) (bus Bus, err error) {
	_, span := app.SpanWithContext(ctx, "Parse device BUS")
	defer span.End()

	switch strings.ToLower(string(value)) {
	case "pci":
		bus = PCI
	case "bluetooth":
		bus = Bluetooth
	case "usb":
		bus = USB
	default:
		err = fmt.Errorf("invalid device bus: %v", value)
	}

	return
}
