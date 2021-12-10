package monitor

import (
	"sync"
	"time"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/audio"
	"github.com/sadesyllas/go-cctl/app/pubsub"
)

var started = false
var singletonLock sync.Mutex

func Start() {
	defer func() { app.Logger.Fatalf("Audio monitor has stopped\n") }()

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

	for {
		pubsub.Send(pubsub.NewMessage(pubsub.TopicDeviceState, audio.FetchCardsWithDevices()))

		time.Sleep(15 * time.Second)
	}
}
