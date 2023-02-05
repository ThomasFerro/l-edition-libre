package persistency

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptsHistory struct {
	history map[application.ManuscriptID][]events.Event
}

func (manuscripts ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]events.Event, error) {
	return manuscripts.history[manuscriptID], nil
}
func (manuscripts ManuscriptsHistory) Append(manuscriptID application.ManuscriptID, newEvents []events.Event) error {
	persistedEvents, exists := manuscripts.history[manuscriptID]
	if !exists {
		persistedEvents = make([]events.Event, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	manuscripts.history[manuscriptID] = persistedEvents
	return nil
}

func NewManuscriptsHistory() ManuscriptsHistory {
	return ManuscriptsHistory{
		history: make(map[application.ManuscriptID][]events.Event),
	}
}
