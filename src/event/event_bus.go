package event

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type EventBus struct {
	Subscribers map[string]EventSubscriber

	EventChan chan *Event

	stopC chan struct{}
	Lock  sync.Mutex
}

func New() *EventBus {
	bus := &EventBus{
		Subscribers: make(map[string]EventSubscriber),
		EventChan:   make(chan *Event, 1024),
		stopC:       make(chan struct{}, 1),
		Lock:        sync.Mutex{},
	}

	return bus
}

func (bus *EventBus) Start(ctx context.Context) error {
	for {
		select {
		case e := <-bus.EventChan:
			for _, subscriber := range bus.Subscribers {
				if subscriber.InterestIn(e) {
					if err := subscriber.Write(e); err != nil {
						logrus.Debugf("write event e %s to %s got error: %s", e, subscriber, err)
					} else {
						logrus.Debugf("write event e %s to %s", e, subscriber)
					}
				} else {
					logrus.Debugf("subscriber %s have no interest in %s", subscriber, e)
				}
			}

		case <-bus.stopC:
			return nil

		case <-ctx.Done():
			return nil
		}
	}
}

func (bus *EventBus) Stop() {
	bus.stopC <- struct{}{}
}
