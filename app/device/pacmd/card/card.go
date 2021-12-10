package card

import (
	"github.com/sadesyllas/go-cctl/app/device/bus"
	"github.com/sadesyllas/go-cctl/app/device/formfactor"
)

type Card struct {
	Index         uint64                `json:"index"`
	Name          string                `json:"name"`
	Driver        string                `json:"driver"`
	Description   string                `json:"description"`
	Profiles      []CardProfile         `json:"profiles"`
	ActiveProfile CardProfile           `json:"activeProfile"`
	SourceIds     []uint64              `json:"sourceIds"`
	SinkIds       []uint64              `json:"sinkIds"`
	FormFactor    formfactor.FormFactor `json:"formFactor"`
	Bus           bus.Bus               `json:"bus"`
}
