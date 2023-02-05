package application

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
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
