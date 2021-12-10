package carddevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

type BluetoothProtocol uint64

const (
	HeadsetHeadUnit BluetoothProtocol = iota + 1
	A2DPSink
)

type ParseableBluetoothProtocol string

func (value ParseableBluetoothProtocol) Parse(ctx context.Context) (bluetoothProtocol BluetoothProtocol, err error) {
	_, span := app.SpanWithContext(ctx, "Parse Device Bluetooth Protocol")
	defer span.End()

	switch strings.ToLower(string(value)) {
	case "headset_head_unit":
		bluetoothProtocol = HeadsetHeadUnit
	case "a2dp_sink":
		bluetoothProtocol = A2DPSink
	default:
		err = fmt.Errorf("Invalid device bluetooth protocol: %v\n", value)
	}

	return
}
