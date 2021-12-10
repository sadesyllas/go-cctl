package carddevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

type A2DPCodec uint64

const (
	SBC A2DPCodec = iota + 1
	AAC
	AptX
)

type ParseableA2DPCodec string

func (value ParseableA2DPCodec) Parse(ctx context.Context) (a2dpCodec A2DPCodec, err error) {
	_, span := app.SpanWithContext(ctx, "Parse Device Bluetooth Protocol")
	defer span.End()

	switch strings.ToLower(string(value)) {
	case "sbc":
		a2dpCodec = SBC
	case "aac":
		a2dpCodec = AAC
	case "aptx":
		a2dpCodec = AptX
	default:
		err = fmt.Errorf("Invalid device A2DP codec: %v\n", value)
	}

	return
}
