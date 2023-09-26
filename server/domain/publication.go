package domain

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type Publication struct {
	Status PublicationStatus
}

type PublicationStatus string

const (
	Available PublicationStatus = "Available"
)

func MakePublicationAvailable() ([]events.Event, DomainError) {
	return []events.Event{
		events.PublicationMadeAvailable{},
	}, nil
}

func (p Publication) String() string {
	return fmt.Sprintf("Publication{Status %v}", p.Status)
}
