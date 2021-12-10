package carddevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

type DeviceState uint64

const (
	Running DeviceState = iota + 1
	Idle
	Suspended
)

type ParseableDeviceState string

func (value ParseableDeviceState) Parse(ctx context.Context) (deviceState DeviceState, err error) {
	_, span := app.SpanWithContext(ctx, "Parse Device State")
	defer span.End()

	switch strings.ToLower(string(value)) {
	case "running":
		deviceState = Running
	case "idle":
		deviceState = Idle
	case "suspended":
		deviceState = Suspended
	default:
		err = fmt.Errorf("Invalid device state: %v\n", value)
	}

	return
}
