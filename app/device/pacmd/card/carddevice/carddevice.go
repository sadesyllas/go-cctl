package carddevice

import (
	"github.com/sadesyllas/go-cctl/app/device/bus"
	"github.com/sadesyllas/go-cctl/app/device/formfactor"
)

type CardDevice struct {
	Index             uint64                `json:"index"`
	Name              string                `json:"name"`
	Driver            string                `json:"driver"`
	State             DeviceState           `json:"state"`
	IsDefault         bool                  `json:"isDefault"`
	Volume            float64               `json:"volume"`
	IsMuted           bool                  `json:"isMuted"`
	CardIndex         uint64                `json:"cardIndex"`
	Description       string                `json:"description"`
	BluetoothProtocol BluetoothProtocol     `json:"bluetoothProtocol"`
	A2DPCodec         A2DPCodec             `json:"a2dpCodec"`
	FormFactor        formfactor.FormFactor `json:"formFactor"`
	Bus               bus.Bus               `json:"bus"`
}
