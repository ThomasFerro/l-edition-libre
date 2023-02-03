package queries

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptStatus struct {
	// TODO: L'ID devrait être géré un cran au-dessus, et on ne reçoit que les event de cet ID
	events.ManuscriptID
}

func GetManuscriptStatus(history []events.Event, query ManuscriptStatus) (domain.Status, error) {
	return domain.Rehydrate(history).Status, nil
}
