package card

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/bus"
	"github.com/sadesyllas/go-cctl/app/device/formfactor"
	"github.com/sadesyllas/go-cctl/app/util"
)

func Parse(text string, ctx context.Context) []*Card {
	ctx, span := app.SpanWithContext(ctx, "Parse Cards")
	defer span.End()

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	re := regexp.MustCompile(`(?P<key>(?:\*\s*)?[^:=]+?)\s*[:=]\s*(?P<value>.+$)?`)
	cards := []*Card{}
	card := (*Card)(nil)
	inProfiles := false
	inSinks := false
	inSources := false

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
		case "index":
			if card != nil {
				cards = append(cards, card)
			}

			card = new(Card)

			value := match[captures["value"]]

			card.Index, _ = strconv.ParseUint(value, 10, 0)
		case "name", "driver", "device.description":
			value := util.UnquoteParsedStringValue(match[captures["value"]])

			switch key {
			case "name":
				card.Name = value
			case "driver":
				card.Driver = value
			case "device.description":
				card.Description = value
			}
		case "profiles":
			inProfiles = true
		case "active profile":
			inProfiles = false

			if card.Bus != bus.Bluetooth {
				continue
			}

			value := util.UnquoteParsedStringValue(match[captures["value"]])
			profile, _ := ParseableProfile(value).Parse()

			card.ActiveProfile = profile
		case "sinks":
			inSinks = true
		case "sources":
			inSinks = false
			inSources = true
		case "ports":
			inSources = false
		case "device.bus":
			value := util.UnquoteParsedStringValue(match[captures["value"]])
			bus, _ := bus.ParseableBus(value).Parse(ctx)

			card.Bus = bus
		case "device.form_factor":
			value := util.UnquoteParsedStringValue(match[captures["value"]])
			formFactor, _ := formfactor.ParseableFormFactor(value).Parse(ctx)

			card.FormFactor = formFactor
		default:
			if inProfiles && card.Bus == bus.Bluetooth {
				profile, _ := ParseableProfile(key).Parse()
				card.Profiles = append(card.Profiles, profile)
			}

			if inSinks || inSources {
				value, _ := strconv.ParseUint(strings.Split(key, "#")[1], 10, 0)

				if inSinks {
					card.SinkIds = append(card.SinkIds, value)
				} else if inSources {
					card.SourceIds = append(card.SourceIds, value)
				}
			}
		}
	}

	if card != nil {
		cards = append(cards, card)
	}

	return cards
}
