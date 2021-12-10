package card

import (
	"fmt"
	"strings"
)

type CardProfile uint64

const (
	HeadsetHeadUnit CardProfile = iota + 1
	A2DPSinkSBC
	A2DPSinkAAC
	A2DPSinkAptX
	A2DPSinkAptXHD
	A2DPSinkLDAC
	Off
)

type ParseableProfile string

func (value ParseableProfile) Parse() (cardProfile CardProfile, err error) {
	switch strings.ToLower(string(value)) {
	case "headset_head_unit":
		cardProfile = HeadsetHeadUnit
	case "a2dp_sink_sbc":
		cardProfile = A2DPSinkSBC
	case "a2dp_sink_aac":
		cardProfile = A2DPSinkAAC
	case "a2dp_sink_aptx":
		cardProfile = A2DPSinkAptX
	case "a2dp_sink_aptx_hd":
		cardProfile = A2DPSinkAptXHD
	case "a2dp_sink_ldac":
		cardProfile = A2DPSinkLDAC
	case "off":
		cardProfile = Off
	default:
		err = fmt.Errorf("invalid card profile: %v\n", value)
	}

	return
}

func (value CardProfile) String() string {
	switch value {
	case HeadsetHeadUnit:
		return "headset_head_unit"
	case A2DPSinkSBC:
		return "a2dp_sink_sbc"
	case A2DPSinkAAC:
		return "a2dp_sink_aac"
	case A2DPSinkAptX:
		return "a2dp_sink_aptx"
	case A2DPSinkAptXHD:
		return "a2dp_sink_aptx_hd"
	case A2DPSinkLDAC:
		return "a2dp_sink_ldac"
	}

	panic("unreachable")
}
