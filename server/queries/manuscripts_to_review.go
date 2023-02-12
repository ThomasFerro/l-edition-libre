package queries

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptsToReview struct{}

func GetManuscriptsToReview(historyForManuscripts [][]events.Event, query ManuscriptsToReview) ([]domain.Manuscript, error) {
	manuscripts := make([]domain.Manuscript, 0)
	for _, historyForManuscript := range historyForManuscripts {
		manuscript := domain.Rehydrate(historyForManuscript)
		if manuscript.Status == domain.PendingReview {
			manuscripts = append(manuscripts, manuscript)
		}
	}
	return manuscripts, nil
}
