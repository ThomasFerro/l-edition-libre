package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type WriterManuscripts struct{}

func HandleWriterManuscripts(ctx context.Context, query Query) (interface{}, error) {
	historyForManuscripts := contexts.FromContext[[][]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptsHistoryContextKey)
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscripts = append(manuscripts, domain.Rehydrate(events.ToEvents(historyForManuscript)))
	}
	return manuscripts, nil
}
