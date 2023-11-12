package application

import (
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type EventContext struct {
	contexts.UserID
}
type ContextualizedEvent struct {
	OriginalEvent events.Event
	Context       EventContext
	ManuscriptID  contexts.ManuscriptID
}

func (e ContextualizedEvent) Event() events.Event {
	return e.OriginalEvent
}
