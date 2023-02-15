package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
)

type WriterManuscripts struct{}

func HandleWriterManuscripts(ctx context.Context, query Query) (interface{}, error) {
	historyForManuscripts := contexts.ManuscriptsHistoryFromContext(ctx)
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscripts = append(manuscripts, domain.Rehydrate(historyForManuscript))
	}
	return manuscripts, nil
}
