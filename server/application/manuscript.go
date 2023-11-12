package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/utils"
)

type Manuscripts interface {
	Persists(contexts.ManuscriptID, domain.Manuscript) error
}

func isTheManuscriptWriter(ctx context.Context) (bool, error) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey{})
	userID := ctx.Value(contexts.UserIDContextKey{}).(contexts.UserID)
	for _, nextEvent := range history {
		_, manuscriptSubmittedEvent := nextEvent.Event().(events.ManuscriptSubmitted)
		if !manuscriptSubmittedEvent {
			continue
		}
		if nextEvent.(ContextualizedEvent).Context.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

type ManuscriptsHistory interface {
	For(contexts.ManuscriptID) ([]ContextualizedEvent, error)
	ForAll() (utils.OrderedMap[contexts.ManuscriptID, []ContextualizedEvent], error)
	ForAllOfUser(contexts.UserID) (utils.OrderedMap[contexts.ManuscriptID, []ContextualizedEvent], error)
	Append(context.Context, []ContextualizedEvent) error
}
