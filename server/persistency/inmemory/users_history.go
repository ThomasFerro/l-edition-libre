package inmemory

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

type UsersHistory struct {
	history map[contexts.UserID][]application.ContextualizedEvent
}

func (users UsersHistory) For(userID contexts.UserID) ([]application.ContextualizedEvent, error) {
	return users.history[userID], nil
}

func (users UsersHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(contexts.UserID)
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
		history: make(map[contexts.UserID][]application.ContextualizedEvent),
	}
}
