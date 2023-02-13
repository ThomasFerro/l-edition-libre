package queries

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type WriterManuscripts struct{}

func GetWriterManuscripts(historyForManuscripts [][]events.Event, query WriterManuscripts) ([]domain.Manuscript, error) {
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscripts = append(manuscripts, domain.Rehydrate(historyForManuscript))
	}
	return manuscripts, nil
}
