package persistency

import (
	"github.com/ThomasFerro/l-edition-libre/application"
)

type UsersHistory struct {
	history map[application.UserID][]application.ContextualizedEvent
}

func (users UsersHistory) For(userID application.UserID) ([]application.ContextualizedEvent, error) {
	return users.history[userID], nil
}

func (users UsersHistory) Append(userID application.UserID, newEvents []application.ContextualizedEvent) error {
	persistedEvents, exists := users.history[userID]
	if !exists {
		persistedEvents = make([]application.ContextualizedEvent, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	users.history[userID] = persistedEvents
	return nil
}

func NewUsersHistory() application.UsersHistory {
	return UsersHistory{
		history: make(map[application.UserID][]application.ContextualizedEvent),
	}
}
