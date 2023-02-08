package persistency

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptsHistory struct {
	history map[application.UserID]map[application.ManuscriptID][]events.Event
}

func (manuscripts ManuscriptsHistory) forUser(userID application.UserID) map[application.ManuscriptID][]events.Event {
	history, found := manuscripts.history[userID]
	if !found {
		newHistoryForUser := make(map[application.ManuscriptID][]events.Event)
		manuscripts.history[userID] = newHistoryForUser
		return newHistoryForUser
	}
	return history
}

func (manuscripts ManuscriptsHistory) For(userID application.UserID, manuscriptID application.ManuscriptID) ([]events.Event, error) {
	history, found := manuscripts.forUser(userID)[manuscriptID]
	if !found {
		return nil, commands.ManuscriptNotFound{}
	}
	return history, nil
}
func (manuscripts ManuscriptsHistory) Append(userID application.UserID, manuscriptID application.ManuscriptID, newEvents []events.Event) error {
	persistedEvents, exists := manuscripts.forUser(userID)[manuscriptID]
	if !exists {
		persistedEvents = make([]events.Event, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	manuscripts.history[userID][manuscriptID] = persistedEvents
	return nil
}

func NewManuscriptsHistory() ManuscriptsHistory {
	return ManuscriptsHistory{
		history: make(map[application.UserID]map[application.ManuscriptID][]events.Event),
	}
}
