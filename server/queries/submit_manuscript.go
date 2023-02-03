package queries

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptStatus struct {
	// TODO: ID ?
	ManuscriptName string
}

func GetManuscriptStatus(history []events.Event, query ManuscriptStatus) (domain.Status, error) {
	return domain.Rehydrate(history).Status, nil
}
