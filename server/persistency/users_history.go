package persistency

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type UsersHistory struct {
	history map[application.UserID][]events.Event
}

func (users UsersHistory) For(userID application.UserID) ([]events.Event, error) {
	return users.history[userID], nil
}

func (users UsersHistory) Append(userID application.UserID, newEvents []events.Event) error {
	persistedEvents, exists := users.history[userID]
	if !exists {
		persistedEvents = make([]events.Event, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	users.history[userID] = persistedEvents
	return nil
}

func NewUsersHistory() UsersHistory {
	return UsersHistory{
		history: make(map[application.UserID][]events.Event),
	}
}
