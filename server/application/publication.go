package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/google/uuid"
)

type PublicationID uuid.UUID

func (m PublicationID) String() string {
	return uuid.UUID(m).String()
}

func MustParsePublicationID(value string) PublicationID {
	return PublicationID(uuid.MustParse(value))
}

func ParsePublicationID(value string) (PublicationID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return PublicationID{}, err
	}
	return PublicationID(id), nil
}

func NewPublicationID() PublicationID {
	return PublicationID(uuid.New())
}

type Publications interface {
	Persists(PublicationID, domain.Publication) error
}

type PublicationsHistory interface {
	For(PublicationID) ([]ContextualizedEvent, error)
	Append(context.Context, []ContextualizedEvent) error
}
