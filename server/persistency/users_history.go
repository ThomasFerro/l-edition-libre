package persistency

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

type UsersHistory struct {
	history map[application.UserID][]application.ContextualizedEvent
}

func (users UsersHistory) For(userID application.UserID) ([]application.ContextualizedEvent, error) {
	return users.history[userID], nil
}

func (users UsersHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
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
