package queries

import (
	"context"
	"sort"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type WriterManuscripts struct{}

type Manuscript struct {
	Id contexts.ManuscriptID
	domain.Manuscript
}

func HandleWriterManuscripts(ctx context.Context, query Query) (interface{}, error) {
	historyForManuscripts := contexts.FromContext[map[contexts.ManuscriptID][]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptsHistoryContextKey{})
	manuscripts := make([]Manuscript, 0)
	for manuscriptID, historyForManuscript := range historyForManuscripts {

		domainManuscript := domain.RehydrateManuscript(events.ToEvents(historyForManuscript))
		manuscripts = append(manuscripts, Manuscript{
			Id:         manuscriptID,
			Manuscript: domainManuscript,
		})
	}
	sort.Slice(manuscripts, func(i, j int) bool {
		manuscriptA := manuscripts[i]
		manuscriptB := manuscripts[j]
		return manuscriptA.Id.String() < manuscriptB.Id.String()
	})
	return manuscripts, nil
}
