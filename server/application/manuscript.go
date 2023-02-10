package application

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/google/uuid"
)

type ManuscriptID uuid.UUID

func (m ManuscriptID) String() string {
	return uuid.UUID(m).String()
}

func MustParseManuscriptID(value string) ManuscriptID {
	return ManuscriptID(uuid.MustParse(value))
}

func NewManuscriptID() ManuscriptID {
	return ManuscriptID(uuid.New())
}

type Manuscripts interface {
	Persists(ManuscriptID, domain.Manuscript) error
}

func isTheManuscriptWriter(history ManuscriptsHistory, userID UserID, manuscriptID ManuscriptID) (bool, error) {
	forManuscript, err := history.For(manuscriptID)
	if err != nil {
		return false, err
	}
	for _, nextEvent := range forManuscript {
		_, manuscriptSubmittedEvent := nextEvent.Event.(events.ManuscriptSubmitted)
		if !manuscriptSubmittedEvent {
			continue
		}
		if nextEvent.Context.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}
