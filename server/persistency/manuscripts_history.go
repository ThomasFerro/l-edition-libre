package persistency

import (
	"github.com/ThomasFerro/l-edition-libre/application"
)

type ManuscriptsHistory struct {
	history map[application.ManuscriptID][]application.ContextualizedEvent
}

func (manuscripts ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]application.ContextualizedEvent, error) {
	return manuscripts.history[manuscriptID], nil
}

func (manuscripts ManuscriptsHistory) ForAll() (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return manuscripts.history, nil
}

func (manuscripts ManuscriptsHistory) Append(manuscriptID application.ManuscriptID, newEvents []application.ContextualizedEvent) error {
	persistedEvents, exists := manuscripts.history[manuscriptID]
	if !exists {
		persistedEvents = make([]application.ContextualizedEvent, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	manuscripts.history[manuscriptID] = persistedEvents
	return nil
}

func NewManuscriptsHistory() application.ManuscriptsHistory {
	return ManuscriptsHistory{
		history: make(map[application.ManuscriptID][]application.ContextualizedEvent),
	}
}
