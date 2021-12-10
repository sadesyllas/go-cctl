package web

import (
	"time"

	"github.com/sadesyllas/go-cctl/app/device/audio"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
)

type VolumeRequest struct {
	Type   string
	Index  uint64
	Volume float64
}

type MuteRequest struct {
	Type  string
	Index uint64
	Mute  bool
}

type DefaultCardDeviceRequest struct {
	Type  string
	Index uint64
	Name  string
}

type CardProfileRequest struct {
	Index   uint64
	Profile card.CardProfile
}

type CardsWithDevicesResponse struct {
	Cards     []*card.Card             `json:"cards"`
	Sources   []*carddevice.CardDevice `json:"sources"`
	Sinks     []*carddevice.CardDevice `json:"sinks"`
	Timestamp uint64                   `json:"timestamp"`
}

func NewCardsWithDevicesResponse(value *audio.CardsWithDevices) CardsWithDevicesResponse {
	return CardsWithDevicesResponse{
		Cards:     value.Cards,
		Sources:   value.Sources,
		Sinks:     value.Sinks,
		Timestamp: uint64(time.Now().UnixMilli()),
	}
}
