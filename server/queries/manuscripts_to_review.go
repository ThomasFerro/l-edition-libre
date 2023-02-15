package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
)

type ManuscriptsToReview struct{}

func HandleManuscriptsToReview(ctx context.Context, query Query) (interface{}, error) {
	historyForManuscripts := contexts.ManuscriptsHistoryFromContext(ctx)
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscript := domain.Rehydrate(historyForManuscript)
		if manuscript.Status == domain.PendingReview {
			manuscripts = append(manuscripts, manuscript)
		}
	}
	return manuscripts, nil
}
