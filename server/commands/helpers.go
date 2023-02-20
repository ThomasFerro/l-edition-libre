package commands

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func rehydrateFromContext(ctx context.Context) domain.Manuscript {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey{})
	return domain.RehydrateManuscript(events.ToEvents(history))
}
