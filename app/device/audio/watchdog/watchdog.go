package watchdog

import (
	"sync"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/audio"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
	"github.com/sadesyllas/go-cctl/app/pubsub"
)

var started = false
var singletonLock sync.Mutex

func Start() {
	defer func() { app.Logger.Fatalf("Audio watchdog has stopped\n") }()

	doStart := func() bool {
		singletonLock.Lock()
		defer singletonLock.Unlock()
		if started {
			return false
		} else {
			started = true
			return true
		}
	}()

	if !doStart {
		return
	}

	inbound := make(chan pubsub.Message)

	pubsub.Register(pubsub.TopicDeviceState, inbound)

	for msg := range inbound {
		if payload, ok := msg.Payload.(*audio.CardsWithDevices); ok {
			func() {
				ctx, span := app.Span("Watchdog Iteration")
				defer span.End()

				defaultSource, defaultSink := (*carddevice.CardDevice)(nil), (*carddevice.CardDevice)(nil)

				for _, cardDevice := range payload.Sources {
					if cardDevice.IsDefault {
						defaultSource = cardDevice
					}
				}

				for _, cardDevice := range payload.Sinks {
					if cardDevice.IsDefault {
						defaultSink = cardDevice
					}
				}

				if defaultSource != nil {
					audio.MoveAudioClients(carddevice.Source, defaultSource.Index, defaultSource.Name, ctx)
				}

				if defaultSink != nil {
					audio.MoveAudioClients(carddevice.Sink, defaultSink.Index, defaultSink.Name, ctx)
				}
			}()
		}
	}
}
