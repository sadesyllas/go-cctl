package formfactor

import (
	"context"
	"fmt"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

type FormFactor uint64

const (
	Internal FormFactor = iota + 1
	Headphones
	Webcam
	Headset
)

type ParseableFormFactor string

func (value ParseableFormFactor) Parse(ctx context.Context) (formFactor FormFactor, err error) {
	_, span := app.SpanWithContext(ctx, "Parse device form factor")
	defer span.End()

	switch strings.ToLower(string(value)) {
	case "internal":
		formFactor = Internal
	case "headphone":
		formFactor = Headphones
	case "webcam":
		formFactor = Webcam
	case "headset":
		formFactor = Headset
	default:
		err = fmt.Errorf("invalid device form factor: %v", value)
	}

	return
}
