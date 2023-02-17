package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptsToReview struct{}

func HandleManuscriptsToReview(ctx context.Context, query Query) (interface{}, error) {
	historyForManuscripts := contexts.FromContext[[][]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptsHistoryContextKey)
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscript := domain.Rehydrate(events.ToEvents(historyForManuscript))
		if manuscript.Status == domain.PendingReview {
			manuscripts = append(manuscripts, manuscript)
		}
	}
	return manuscripts, nil
}
