package types

import (
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
)

type CommandResultCards struct {
	Success bool
	Cards   []*card.Card
}

type CommandResultCardDevices struct {
	Success     bool
	CardDevices []*carddevice.CardDevice
}
