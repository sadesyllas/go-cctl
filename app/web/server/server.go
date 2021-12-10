package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sadesyllas/go-cctl/app"
	"github.com/sadesyllas/go-cctl/app/device/audio"
	"github.com/sadesyllas/go-cctl/app/device/pacmd/card/carddevice"
	"github.com/sadesyllas/go-cctl/app/pubsub"
	"github.com/sadesyllas/go-cctl/app/web"
)

const metricsPath = "/metrics"

var reqCnt = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cctl_metric_total_calls",
		Help: "Total number of calls.",
	},
	[]string{"status", "method", "path"},
)

var reqDurCnt = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "cctl_metric_duration",
	Help: "Duration of all HTTP requests by status code, method and path.",
	Buckets: []float64{
		0.000000001, // 1ns
		0.000000002,
		0.000000005,
		0.00000001, // 10ns
		0.00000002,
		0.00000005,
		0.0000001, // 100ns
		0.0000002,
		0.0000005,
		0.000001, // 1µs
		0.000002,
		0.000005,
		0.00001, // 10µs
		0.00002,
		0.00005,
		0.0001, // 100µs
		0.0002,
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	},
},
	[]string{"status_code", "method", "path"},
)

func Start(port uint16) {
	webApp := fiber.New()

	webApp.Use(handleMetrics)

	webApp.Get("/", handleCORS(func(c *fiber.Ctx) error { return c.Redirect("/audio") }))

	webApp.Get(metricsPath, adaptor.HTTPHandler(promhttp.Handler()))

	webApp.Options("/audio", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Get("/audio", handleCORS(handleAudioRequest))

	webApp.Options("/audio/ws", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Get("/audio/ws", handleCORS(websocket.New(handleWebsocketRequest)))

	webApp.Options("/audio/volume", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Post("/audio/volume", handleCORS(handleVolumeRequest))

	webApp.Options("/audio/mute", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Post("/audio/mute", handleCORS(handleMuteRequest))

	webApp.Options("/audio/default", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Post("/audio/default", handleCORS(handleDefaultCardDeviceRequest))

	webApp.Options("/audio/profile", handleCORS(func(*fiber.Ctx) error { return nil }))
	webApp.Post("/audio/profile", handleCORS(handleCardProfileRequest))

	webApp.Listen(fmt.Sprintf(":%v", port))
}

func handleMetrics(c *fiber.Ctx) error {
	start := time.Now()
	method := c.Method()
	path := c.Path()

	if strings.HasPrefix(path, metricsPath) {
		return c.Next()
	}

	if err := c.Next(); err != nil {
		reqCnt.WithLabelValues(strconv.Itoa(0), method, path).Inc()

		elapsed := float64(time.Since(start).Nanoseconds()) / 1000000000
		reqDurCnt.WithLabelValues(strconv.Itoa(0), method, path).Observe(elapsed)

		return err
	}

	statusCode := strconv.Itoa(c.Response().StatusCode())

	reqCnt.WithLabelValues(statusCode, method, path).Inc()

	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000000
	reqDurCnt.WithLabelValues(statusCode, method, path).Observe(elapsed)

	return nil
}

func handleCORS(handler func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Context().Response.Header.Add("Access-Control-Allow-Origin", "*")
		c.Context().Response.Header.Add("Access-Control-Allow-Headers", "Content-Type")

		return handler(c)
	}
}

func handleAudioRequest(c *fiber.Ctx) error {
	_, span := app.Span("/audio")
	defer span.End()

	cardsWithDevices := emitDeviceState()
	deviceStateJsonBytes, _ := json.Marshal(web.NewCardsWithDevicesResponse(cardsWithDevices))
	deviceStateJson := string(deviceStateJsonBytes)

	c.Context().SetContentType("application/json")

	return c.SendString(deviceStateJson)
}

func handleWebsocketRequest(c *websocket.Conn) {
	remoteAddr := c.RemoteAddr()
	defer func() { app.Logger.Infof("Websocket connection from %v has been closed\n", remoteAddr) }()

	app.Logger.Infof("New websocket connection from %v\n", remoteAddr)

	_, span := app.Span("/audio/ws")
	defer span.End()
	inbound := make(chan pubsub.Message)

	pubsub.Register(pubsub.TopicDeviceState, inbound)

	for {
		msg := <-inbound

		if msg, ok := msg.Payload.(*audio.CardsWithDevices); ok {
			app.Logger.Debug("Sending message down the websocket")

			if err := c.WriteJSON(msg); err != nil {
				break
			}
		}
	}
}

func handleVolumeRequest(c *fiber.Ctx) error {
	ctx, span := app.Span("/audio/volume")
	defer span.End()

	var volumeRequest web.VolumeRequest
	if err := json.Unmarshal(c.Body(), &volumeRequest); err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request")
	}

	cardDeviceType, err := carddevice.ParseableCardDeviceType(volumeRequest.Type).Parse()
	if err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request: invalid card device type")
	}

	audio.SetVolume(cardDeviceType, volumeRequest.Index, volumeRequest.Volume, ctx)

	emitDeviceState()

	c.SendStatus(200)

	return nil
}

func handleMuteRequest(c *fiber.Ctx) error {
	ctx, span := app.Span("/audio/mute")
	defer span.End()

	var muteRequest web.MuteRequest
	if err := json.Unmarshal(c.Body(), &muteRequest); err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request")
	}

	cardDeviceType, err := carddevice.ParseableCardDeviceType(muteRequest.Type).Parse()
	if err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request: invalid card device type")
	}

	audio.ToggleMute(cardDeviceType, muteRequest.Index, muteRequest.Mute, ctx)

	emitDeviceState()

	c.SendStatus(200)

	return nil
}

func handleDefaultCardDeviceRequest(c *fiber.Ctx) error {
	ctx, span := app.Span("/audio/default")
	defer span.End()

	var defaultCardDeviceRequest web.DefaultCardDeviceRequest
	if err := json.Unmarshal(c.Body(), &defaultCardDeviceRequest); err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request")
	}

	cardDeviceType, err := carddevice.ParseableCardDeviceType(defaultCardDeviceRequest.Type).Parse()
	if err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request: invalid card device type")
	}

	audio.SetDefaultCardDevice(cardDeviceType, defaultCardDeviceRequest.Index, ctx)

	audio.MoveAudioClients(cardDeviceType, defaultCardDeviceRequest.Index, defaultCardDeviceRequest.Name, ctx)

	emitDeviceState()

	c.SendStatus(200)

	return nil
}

func handleCardProfileRequest(c *fiber.Ctx) error {
	ctx, span := app.Span("/audio/profile")
	defer span.End()

	var cardProfileRequest web.CardProfileRequest
	if err := json.Unmarshal(c.Body(), &cardProfileRequest); err != nil {
		c.SendStatus(400)

		return fmt.Errorf("bad volume request")
	}

	audio.SetCardProfile(cardProfileRequest.Index, cardProfileRequest.Profile, ctx)

	emitDeviceState()

	c.SendStatus(200)

	return nil
}

func emitDeviceState() *audio.CardsWithDevices {
	cardsWithDevices := audio.FetchCardsWithDevices()

	pubsub.Send(pubsub.NewMessage(pubsub.TopicDeviceState, cardsWithDevices))

	return cardsWithDevices
}
