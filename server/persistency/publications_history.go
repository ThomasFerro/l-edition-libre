package persistency

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
)

type PublicationsHistory struct {
	history map[application.PublicationID][]application.ContextualizedEvent
}

func (publications PublicationsHistory) For(publicationID application.PublicationID) ([]application.ContextualizedEvent, error) {
	return publications.history[publicationID], nil
}

func (publications PublicationsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	publicationID := ctx.Value(contexts.PublicationIDContextKey{}).(application.PublicationID)
	persistedEvents, exists := publications.history[publicationID]
	if !exists {
		persistedEvents = make([]application.ContextualizedEvent, 0)
	}
	persistedEvents = append(persistedEvents, newEvents...)
	publications.history[publicationID] = persistedEvents
	return nil
}

func NewPublicationsHistory() application.PublicationsHistory {
	return PublicationsHistory{
		history: make(map[application.PublicationID][]application.ContextualizedEvent),
	}
}
