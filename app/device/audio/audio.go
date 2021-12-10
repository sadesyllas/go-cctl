package audio

import (
	"context"
	"fmt"
	"math"
	"os/exec"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/audioclient"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
	"github.com/sadesyllas/go-cctl/app/internal/types"
	"go.opentelemetry.io/otel/attribute"
)

type CardsWithDevices struct {
	Cards   []*card.Card             `json:"cards"`
	Sources []*carddevice.CardDevice `json:"sources"`
	Sinks   []*carddevice.CardDevice `json:"sinks"`
}

func FetchCardsWithDevices() *CardsWithDevices {
	ctx, span := app.Span("FetchDevices")
	defer span.End()

	cardsCh := make(chan types.CommandResultCards)
	go fetchCards(cardsCh, ctx)

	sourcesCh := make(chan types.CommandResultCardDevices)
	go fetchCardDevices(carddevice.Source, sourcesCh, ctx)

	sinksCh := make(chan types.CommandResultCardDevices)
	go fetchCardDevices(carddevice.Sink, sinksCh, ctx)

	resultCards := <-cardsCh
	resultSources := <-sourcesCh
	resultSinks := <-sinksCh

	result := new(CardsWithDevices)

	if resultCards.Success {
		result.Cards = resultCards.Cards
	}

	if resultSources.Success {
		result.Sources = resultSources.CardDevices
	}

	if resultSinks.Success {
		result.Sinks = resultSinks.CardDevices
	}

	return result
}

func SetVolume(t carddevice.CardDeviceType, index uint64, volumePercentage float64, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "SetVolume")
	span.SetAttributes(
		attribute.Int64("type", int64(t)),
		attribute.Int64("index", int64(index)),
		attribute.Float64("volumePercentage", volumePercentage))
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "set-source-volume"
	} else {
		arg = "set-sink-volume"
	}

	volumePercentage = math.Min(volumePercentage, 100)
	volume := fmt.Sprint(uint64(math.Round(math.Round((volumePercentage*65535/100)*10) / 10)))

	out, err := exec.Command("pacmd", arg, fmt.Sprint(index), volume).CombinedOutput()

	if err != nil {
		app.Logger.Errorf("Could not set the volume of %v index %v to %v", t, index, volumePercentage)
	}

	app.Logger.Debugf("pacmd out: %v", string(out))
}

func ToggleMute(t carddevice.CardDeviceType, index uint64, mute bool, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "ToggleMute")
	span.SetAttributes(
		attribute.Int64("type", int64(t)),
		attribute.Int64("index", int64(index)),
		attribute.Bool("mute", mute))
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "set-source-mute"
	} else {
		arg = "set-sink-mute"
	}

	var muteValue string
	if mute {
		muteValue = "1"
	} else {
		muteValue = "0"
	}

	out, err := exec.Command("pacmd", arg, fmt.Sprint(index), muteValue).CombinedOutput()

	if err != nil {
		app.Logger.Errorf("Could not set mute of %v index %v to mute status %v", t, index, muteValue)
	}

	app.Logger.Debugf("pacmd out: %v", string(out))
}

func SetDefaultCardDevice(t carddevice.CardDeviceType, index uint64, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "SetDefaultCardDevice")
	span.SetAttributes(
		attribute.Int64("type", int64(t)),
		attribute.Int64("index", int64(index)))
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "set-default-source"
	} else {
		arg = "set-default-sink"
	}

	out, err := exec.Command("pacmd", arg, fmt.Sprint(index)).CombinedOutput()

	if err != nil {
		app.Logger.Errorf("Could not set the %v index %v as the default %v", t, index, t)
	}

	app.Logger.Debugf("pacmd out: %v", string(out))
}

func SetCardProfile(index uint64, profile card.CardProfile, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "SetCardProfile")
	span.SetAttributes(
		attribute.Int64("index", int64(index)),
		attribute.Int64("profile", int64(index)))
	defer span.End()

	out, err := exec.Command("pacmd", "set-card-profile", fmt.Sprint(index), fmt.Sprint(profile)).CombinedOutput()

	if err != nil {
		app.Logger.Errorf("Could not set the card index %v to profile %v", index, profile)
	}

	app.Logger.Debugf("pacmd out: %v", string(out))
}

func MoveAudioClients(t carddevice.CardDeviceType, index uint64, name string, ctx context.Context) {
	ctx, span := app.SpanWithContext(ctx, "MoveAudioClients")
	span.SetAttributes(
		attribute.Int64("type", int64(t)),
		attribute.Int64("index", int64(index)),
		attribute.String("name", name))
	defer span.End()

	audioClients := fetchAudioClients(t, ctx)

	for _, audioClient := range audioClients {
		if audioClient.CardDeviceIndex != index {
			connectAudioClientToCardDevice(*audioClient, t, name, ctx)

			app.Logger.Infof("Moved audio client index %v to default %v %v",
				audioClient.Index, t, name)
		}
	}
}

func fetchCards(ch chan<- types.CommandResultCards, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "fetchCards")
	defer span.End()

	out, err := exec.Command("pacmd", "list-cards").CombinedOutput()

	ch <- types.CommandResultCards{
		Success: err == nil,
		Cards:   card.Parse(string(out), ctx)}
}

func fetchCardDevices(t carddevice.CardDeviceType, ch chan<- types.CommandResultCardDevices, ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "fetchCardDevices")
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "list-sources"
	} else {
		arg = "list-sinks"
	}

	out, err := exec.Command("pacmd", arg).CombinedOutput()

	ch <- types.CommandResultCardDevices{
		Success:     err == nil,
		CardDevices: carddevice.Parse(string(out), ctx)}
}

func fetchAudioClients(t carddevice.CardDeviceType, ctx context.Context) []*audioclient.AudioClient {
	ctx, span := app.SpanWithContext(ctx, "fetchClientIndexes")
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "list-source-outputs"
	} else {
		arg = "list-sink-inputs"
	}

	out, err := exec.Command("pacmd", arg).CombinedOutput()
	if err != nil {
		app.Logger.Fatalf("Failed to get the audio client indexes\n")
	}

	return audioclient.Parse(string(out), ctx)
}

func connectAudioClientToCardDevice(
	audioClient audioclient.AudioClient,
	t carddevice.CardDeviceType,
	cardDeviceName string,
	ctx context.Context) {
	_, span := app.SpanWithContext(ctx, "connectAudioClientToCardDevice")
	defer span.End()

	var arg string
	if t == carddevice.Source {
		arg = "move-source-output"
	} else {
		arg = "move-sink-input"
	}

	_, err := exec.Command("pacmd", arg, fmt.Sprint(audioClient.Index), cardDeviceName).CombinedOutput()
	if err != nil {
		app.Logger.Errorf("Could not set client index %v to %v %v", audioClient.Index, t, cardDeviceName)
	}
}
