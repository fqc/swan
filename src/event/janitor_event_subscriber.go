package event

import (
	"encoding/json"
	"sync"

	"github.com/Dataman-Cloud/swan/src/types"
	"github.com/Sirupsen/logrus"

	"github.com/Dataman-Cloud/swan-janitor/src"
)

type JanitorSubscriber struct {
	Key          string
	acceptors    map[string]types.JanitorAcceptor
	acceptorLock sync.RWMutex
}

func NewJanitorSubscriber() *JanitorSubscriber {
	janitorSubscriber := &JanitorSubscriber{
		Key:       "janitor",
		acceptors: make(map[string]types.JanitorAcceptor),
	}
	return janitorSubscriber
}

func (js *JanitorSubscriber) Subscribe(bus *EventBus) error {
	bus.Lock.Lock()
	defer bus.Lock.Unlock()

	bus.Subscribers[js.Key] = js
	return nil
}

func (js *JanitorSubscriber) Unsubscribe(bus *EventBus) error {
	bus.Lock.Lock()
	defer bus.Lock.Unlock()

	delete(bus.Subscribers, js.Key)
	return nil
}

func (js *JanitorSubscriber) AddAcceptor(acceptor types.JanitorAcceptor) {
	js.acceptorLock.Lock()
	js.acceptors[acceptor.ID] = acceptor
	js.acceptorLock.Unlock()
}

func (js *JanitorSubscriber) RemoveAcceptor(ID string) {
	js.acceptorLock.Lock()
	delete(js.acceptors, ID)
	js.acceptorLock.Unlock()
}

func (js *JanitorSubscriber) Write(e *Event) error {
	janitorEvent, err := BuildJanitorEvent(e)
	if err != nil {
		return err
	}

	go js.pushJanitorEvent(janitorEvent)

	return nil
}

func (js *JanitorSubscriber) InterestIn(e *Event) bool {
	if e.AppMode != "replicates" {
		return false
	}

	if e.Type == EventTypeTaskHealthy {
		return true
	}

	if e.Type == EventTypeTaskUnhealthy {
		return true
	}

	return false
}

func (js *JanitorSubscriber) pushJanitorEvent(event *janitor.TargetChangeEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		logrus.Infof("marshal janitor event got error: %s", err.Error())
		return
	}

	js.acceptorLock.RLock()
	for _, acceptor := range js.acceptors {
		if err := SendEventByHttp(acceptor.RemoteAddr, "POST", data); err != nil {
			logrus.Infof("send janitor event by http to %s got error: %s", acceptor.RemoteAddr, err.Error())
		}
	}
	js.acceptorLock.RUnlock()
}
