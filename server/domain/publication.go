package domain

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
	"golang.org/x/exp/slog"
)

type Publication struct {
	Status PublicationStatus
}

type PublicationStatus string

const (
	Available PublicationStatus = "Available"
)

func (p Publication) applyPublicationMadeAvailable(e events.PublicationMadeAvailable) Publication {
	p.Status = Available
	return p
}

func (p Publication) String() string {
	return fmt.Sprintf("Publication{Status %v}", p.Status)
}

func RehydratePublication(history []events.Event) Publication {
	publication := Publication{}

	for _, nextEvent := range history {
		switch typedEvent := nextEvent.(type) {
		case events.PublicationMadeAvailable:
			publication = publication.applyPublicationMadeAvailable(typedEvent)
		default:
			slog.Warn("unknown publication event", "event", typedEvent)
		}
	}

	return publication
}
