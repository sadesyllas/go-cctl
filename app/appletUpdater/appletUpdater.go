package appletUpdater

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/audio"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
	"github.com/sadesyllas/go-cctl/app/pubsub"
)

var started = false
var singletonLock sync.Mutex

func Start() {
	defer func() { app.Logger.Fatalf("Applet updater has stopped\n") }()

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

	defaultSource := (*carddevice.CardDevice)(nil)

	for msg := range inbound {
		if payload, ok := msg.Payload.(*audio.CardsWithDevices); ok {
			func() {
				_, span := app.Span("Applet Updater Iteration")
				defer span.End()

				newDefaultSource := (*carddevice.CardDevice)(nil)

				for _, cardDevice := range payload.Sources {
					if cardDevice.IsDefault {
						newDefaultSource = cardDevice
					}
				}

				if defaultSource == nil ||
					(newDefaultSource.Index != defaultSource.Index ||
						newDefaultSource.Volume != defaultSource.Volume ||
						newDefaultSource.IsMuted != defaultSource.IsMuted) {
					defaultSource = newDefaultSource

					var volumeIcon string
					if newDefaultSource.IsMuted {
						volumeIcon = "microphone-sensitivity-muted-symbolic"
					} else if newDefaultSource.Volume < 25.0 {
						volumeIcon = "microphone-sensitivity-low-symbolic"
					} else if newDefaultSource.Volume > 75.0 {
						volumeIcon = "microphone-sensitivity-high-symbolic"
					} else {
						volumeIcon = "microphone-sensitivity-medium-symbolic"
					}

					homeEnvVar, _ := os.LookupEnv("HOME")
					appFilePathGlob := path.Join(homeEnvVar, ".config", "xfce4", "panel", "**", "*.desktop")
					appFilePaths, _ := filepath.Glob(appFilePathGlob)

					var appletFilePath string
					for _, appFilePath := range appFilePaths {
						content, _ := ioutil.ReadFile(appFilePath)
						if strings.Contains(string(content), "Name=toggle_microphone") {
							appletFilePath = appFilePath

							break
						}
					}

					if appletFilePath != "" {
						cmd := exec.Command("sed", "-i", fmt.Sprintf("s/Icon=.*/Icon=%v/", volumeIcon), appletFilePath)
						cmd.CombinedOutput()

						if !cmd.ProcessState.Success() {
							app.Logger.Errorf("Could not set the applet icon to %v\n", volumeIcon)
						}
					}

					cmd := exec.Command("notify-send", "-t", "1", "-i", volumeIcon, fmt.Sprint(newDefaultSource.Volume))
					cmd.CombinedOutput()

					if !cmd.ProcessState.Success() {
						app.Logger.Errorf("Could not notify about new default source state\n")
					}
				}
			}()
		}
	}
}
