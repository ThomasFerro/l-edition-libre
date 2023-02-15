package queries

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
)

type ManuscriptStatus struct{}

func HandleManuscriptStatus(ctx context.Context, query Query) (interface{}, error) {
	history := contexts.ManuscriptHistoryFromContext(ctx)
	return domain.Rehydrate(history).Status, nil
}
