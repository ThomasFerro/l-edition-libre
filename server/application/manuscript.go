package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/utils"
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

func isTheManuscriptWriter(ctx context.Context) (bool, error) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey{})
	userID := ctx.Value(contexts.UserIDContextKey{}).(UserID)
	for _, nextEvent := range history {
		_, manuscriptSubmittedEvent := nextEvent.Event().(events.ManuscriptSubmitted)
		if !manuscriptSubmittedEvent {
			continue
		}
		if nextEvent.(ContextualizedEvent).Context.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

type ManuscriptsHistory interface {
	For(ManuscriptID) ([]ContextualizedEvent, error)
	ForAll() (utils.OrderedMap[ManuscriptID, []ContextualizedEvent], error)
	ForAllOfUser(UserID) (utils.OrderedMap[ManuscriptID, []ContextualizedEvent], error)
	Append(context.Context, []ContextualizedEvent) error
}
