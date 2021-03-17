package evb

import "github.com/asaskevich/EventBus"

type Event struct {
	Action string
	Data   interface{}
}

var bus = EventBus.New()

func Subscribe(actorID string, cb func(e Event)) {
	bus.Subscribe(actorID, cb)
}

func Publish(actorID string, e Event) {
	bus.Publish(actorID, e)
}
