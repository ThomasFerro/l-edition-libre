package contexts

import "github.com/google/uuid"

type UserID string
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
