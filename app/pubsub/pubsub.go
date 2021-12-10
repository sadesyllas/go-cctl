package pubsub

import (
	"math"
	"sync"
	"time"

	"github.com/sadesyllas/go-cctl/app"
)

type Topic uint64

const (
	TopicDeviceState Topic = iota + 1
)

type Message struct {
	Topic   Topic
	Payload interface{}
}

func NewMessage(topic Topic, payload interface{}) Message {
	return Message{
		Topic:   topic,
		Payload: payload,
	}
}

var inbound = make(chan Message)
var subscriptions = make(map[Topic][]chan<- Message)
var singletonLock sync.Mutex
var started = false

func Start() {
	defer func() { app.Logger.Fatal("PubSub has stopped") }()

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

	for msg := range inbound {
		func() {
			singletonLock.Lock()
			defer singletonLock.Unlock()
			if ch, ok := subscriptions[msg.Topic]; ok {
				for i, c := range ch {
					timer := time.NewTimer(time.Second)
					timerExpired := false
					done := false

					for {
						if done {
							break
						}

						select {
						case c <- msg:
							done = true
						default:
							select {
							case <-timer.C:
								timerExpired = true
								done = true
							default:
								time.Sleep(10 * time.Millisecond)
							}
						}
					}

					if timerExpired {
						subscriptions[msg.Topic] = append(ch[:i], ch[int(math.Min(float64(i+1), float64(len(ch)))):]...)

						app.Logger.Info("PubSub evicted timed out receiving channel")
					}
				}
			}
		}()
	}
}

func Register(topic Topic, ch chan<- Message) {
	singletonLock.Lock()
	defer singletonLock.Unlock()

	if _, ok := subscriptions[topic]; !ok {
		subscriptions[topic] = make([]chan<- Message, 0)
	}

	subscriptions[topic] = append(subscriptions[topic], ch)
}

func Send(msg Message) {
	inbound <- msg
}
