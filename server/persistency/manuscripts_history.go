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

func (manuscripts ManuscriptsHistory) ForAllOfUser(userID application.UserID) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	userHistory := map[application.ManuscriptID][]application.ContextualizedEvent{}
	for manuscriptID, manuscriptHistory := range manuscripts.history {
		userManuscriptHistory := []application.ContextualizedEvent{}
		for _, manuscriptEvent := range manuscriptHistory {
			if manuscriptEvent.Context.UserID.String() == userID.String() {
				userManuscriptHistory = append(userManuscriptHistory, manuscriptEvent)
			}
		}
		if len(userManuscriptHistory) != 0 {
			userHistory[manuscriptID] = userManuscriptHistory
		}
	}
	return userHistory, nil
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
