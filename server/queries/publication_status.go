package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type PublicationStatus struct{}

func HandlePublicationStatus(ctx context.Context, query Query) (interface{}, error) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedPublicationHistoryContextKey{})
	return domain.RehydratePublication(events.ToEvents(history)).Status, nil
}
