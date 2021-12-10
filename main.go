package main

import (
	"os"
	"sync"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/appletUpdater"
	"github.com/sadesyllas/go-cctl/app/device/audio/monitor"
	"github.com/sadesyllas/go-cctl/app/device/audio/watchdog"
	"github.com/sadesyllas/go-cctl/app/pubsub"
	"github.com/sadesyllas/go-cctl/app/web/server"
	"github.com/spf13/pflag"
)

func main() {
	port := pflag.Uint16P("port", "p", 0, "The web server port")
	pflag.Parse()

	if *port == 0 {
		pflag.Usage()

		os.Exit(1)
	}

	app.SetupLogging()

	stopTracing := app.SetupTracing()
	defer stopTracing()

	var wait sync.WaitGroup
	wait.Add(1)

	go pubsub.Start()
	go monitor.Start()
	go watchdog.Start()
	go appletUpdater.Start()
	go server.Start(*port)

	// tc := make(chan pubsub.Message)
	// pubsub.Register(pubsub.TopicDeviceState, tc)

	// var msg pubsub.Message
	// for {
	// 	msg = <-tc
	// 	j, _ := json.Marshal(msg.Payload.(*audio.CardsWithDevices))
	// 	app.Logger.Debugf("[PUBSUB MSG] %v\n", string(j))
	// }

	wait.Wait()
}
