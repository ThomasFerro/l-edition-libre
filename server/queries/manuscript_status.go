package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptState struct{}

func HandleManuscriptState(ctx context.Context, query Query) (interface{}, error) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey{})
	return domain.RehydrateManuscript(events.ToEvents(history)), nil
}
