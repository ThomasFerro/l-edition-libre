package application

import (
	"context"

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

func ParseManuscriptID(value string) (ManuscriptID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return ManuscriptID{}, err
	}
	return ManuscriptID(id), nil
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

func toEventsByManuscript(toMap map[ManuscriptID][]ContextualizedEvent) [][]events.Event {
	returned := make([][]events.Event, 0)
	for _, nextManuscript := range toMap {
		mappedEvents := make([]events.Event, 0)
		for _, nextEvent := range nextManuscript {
			mappedEvents = append(mappedEvents, nextEvent.Event)
		}
		returned = append(returned, mappedEvents)
	}
	return returned
}

type ManuscriptsHistory interface {
	For(ManuscriptID) ([]ContextualizedEvent, error)
	ForAll() (map[ManuscriptID][]ContextualizedEvent, error)
	ForAllOfUser(UserID) (map[ManuscriptID][]ContextualizedEvent, error)
	Append(context.Context, []ContextualizedEvent) error
}
