package events

type Event interface{}

type UserEvent interface {
	UserEventName() string
}

type ManuscriptEvent interface {
	ManuscriptEventName() string
}

type DecoratedEvent interface {
	Event() Event
}

func ToEvents(decorated []DecoratedEvent) []Event {
	returned := []Event{}
	for _, nextEvent := range decorated {
		returned = append(returned, ToEvent(nextEvent))
	}
	return returned
}

func ToEvent(decorated DecoratedEvent) Event {
	return decorated.Event()
}
