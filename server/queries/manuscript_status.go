package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptStatus struct{}

func HandleManuscriptStatus(ctx context.Context, query Query) (interface{}, error) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey)
	return domain.Rehydrate(events.ToEvents(history)).Status, nil
}
