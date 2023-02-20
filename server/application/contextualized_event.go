package application

import "github.com/ThomasFerro/l-edition-libre/events"

type EventContext struct {
	UserID
}
type ContextualizedEvent struct {
	OriginalEvent events.Event
	Context       EventContext
	ManuscriptID  ManuscriptID
}

func (e ContextualizedEvent) Event() events.Event {
	return e.OriginalEvent
}
