package carddevice

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/bus"
	"github.com/sadesyllas/go-cctl/app/device/formfactor"
	"github.com/sadesyllas/go-cctl/app/util"
)

func Parse(text string, ctx context.Context) []*CardDevice {
	ctx, span := app.SpanWithContext(ctx, "Parse Card Devices")
	defer span.End()

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	re := regexp.MustCompile(`(?P<key>(?:\*\s*)?[^:=]+?)\s*[:=]\s*(?P<value>.+$)`)
	cardDevices := []*CardDevice{}
	cardDevice := (*CardDevice)(nil)

	var captures = make(map[string]int)
	for i, name := range re.SubexpNames() {
		if name != "" {
			captures[name] = i
		}
	}

	for i := 0; i < len(lines); i++ {
		lines[i] = strings.TrimSpace(lines[i])
		match := re.FindStringSubmatch(lines[i])

		if len(match) == 0 {
			continue
		}

		key := match[captures["key"]]

		switch key {
		case "* index", "index":
			if cardDevice != nil {
				cardDevices = append(cardDevices, cardDevice)
			}

			cardDevice = new(CardDevice)

			value := match[captures["value"]]

			cardDevice.Index, _ = strconv.ParseUint(value, 10, 0)

			if strings.HasPrefix(key, "*") {
				cardDevice.IsDefault = true
			}
		case "name", "driver", "device.description":
			if cardDevice == nil {
				break
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])

			switch key {
			case "name":
				cardDevice.Name = value
			case "driver":
				cardDevice.Driver = value
			case "device.description":
				cardDevice.Description = value
			}
		case "muted":
			if cardDevice == nil {
				break
			}

			cardDevice.IsMuted = match[captures["value"]] == "yes"
		case "card":
			if cardDevice == nil {
				break
			}

			value, _ := strconv.ParseUint(strings.Split(match[captures["value"]], " ")[0], 10, 0)

			cardDevice.CardIndex = value
		case "device.form_factor":
			if cardDevice == nil {
				break
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])
			formFactor, _ := formfactor.ParseableFormFactor(value).Parse(ctx)

			cardDevice.FormFactor = formFactor
		case "state":
			if cardDevice == nil {
				break
			}

			value := match[captures["value"]]
			deviceState, _ := ParseableDeviceState(value).Parse(ctx)

			cardDevice.State = deviceState
		case "volume":
			if cardDevice == nil {
				break
			}

			value := match[captures["value"]]
			re := regexp.MustCompile(`^[^:]+:\s*(?P<volume>[0-9]+).*`)
			match := re.FindStringSubmatch(value)

			var captures = make(map[string]int)
			for i, name := range re.SubexpNames() {
				if name != "" {
					captures[name] = i
				}
			}

			volume, _ := strconv.ParseFloat(match[captures["volume"]], 64)
			volume = math.Round((volume / 65535.0) * 100.0)

			cardDevice.Volume = volume
		case "bluetooth.protocol":
			if cardDevice == nil {
				break
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])
			bluetoothProtocol, _ := ParseableBluetoothProtocol(value).Parse(ctx)

			cardDevice.BluetoothProtocol = bluetoothProtocol
		case "bluetooth.a2dp_codec":
			if cardDevice == nil {
				break
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])
			a2dpCodec, _ := ParseableA2DPCodec(value).Parse(ctx)

			cardDevice.A2DPCodec = a2dpCodec
		case "device.bus":
			if cardDevice == nil {
				break
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])
			bus, _ := bus.ParseableBus(value).Parse(ctx)

			cardDevice.Bus = bus
		case "monitor_of":
			cardDevice = nil
		}
	}

	if cardDevice != nil {
		cardDevices = append(cardDevices, cardDevice)
	}

	return cardDevices
}
