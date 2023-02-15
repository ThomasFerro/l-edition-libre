package events

type Event interface{}

type UserEvent interface {
	UserEventName() string
}

type ManuscriptEvent interface {
	ManuscriptEventName() string
}
