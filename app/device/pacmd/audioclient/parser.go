package audioclient

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/sadesyllas/go-cctl/app"
)

func Parse(text string, ctx context.Context) []*AudioClient {
	_, span := app.SpanWithContext(ctx, "Parse Audio Clients")
	defer span.End()

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	re := regexp.MustCompile(`(?P<key>(?:\*\s*)?[^:=]+?)\s*[:=]\s*(?P<value>.+$)`)
	audioClients := []*AudioClient{}
	clientIndex, cardDeviceIndex, monitor := uint64(0), uint64(0), false

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
		value := match[captures["value"]]

		switch key {
		case "index":
			value, _ := strconv.ParseUint(value, 10, 0)

			clientIndex = value
		case "source", "sink":
			re := regexp.MustCompile(`\s*(?P<index>[0-9]+)\s*<(?P<name>[^>]+)>`)

			var captures = make(map[string]int)
			for i, name := range re.SubexpNames() {
				if name != "" {
					captures[name] = i
				}
			}

			match := re.FindStringSubmatch(value)

			indexValue, _ := strconv.ParseUint(match[captures["index"]], 10, 0)
			cardDeviceIndex = indexValue

			monitor = strings.HasSuffix(match[captures["name"]], ".monitor")
		case "client":
			if monitor {
				continue
			}

			re := regexp.MustCompile(`\s*(?:[0-9]+)\s*<(?P<name>[^>]+)>`)

			var captures = make(map[string]int)
			for i, name := range re.SubexpNames() {
				if name != "" {
					captures[name] = i
				}
			}

			match := re.FindStringSubmatch(value)

			if match[captures["name"]] == "PulseAudio Volume Control" {
				continue
			}

			audioClients = append(audioClients, &AudioClient{
				Index:           clientIndex,
				CardDeviceIndex: cardDeviceIndex,
			})
		}
	}

	return audioClients
}
