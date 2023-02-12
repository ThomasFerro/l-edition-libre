package application

import (
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/google/uuid"
)

type UserID uuid.UUID

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

func MustParseUserID(value string) UserID {
	return UserID(uuid.MustParse(value))
}

func ParseUserID(value string) (UserID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, err
	}
	return UserID(parsed), nil
}

func NewUserID() UserID {
	return UserID(uuid.New())
}

func userIsAnEditor(history UsersHistory, userID UserID) (bool, error) {
	forUser, err := history.For(userID)
	if err != nil {
		return false, err
	}
	for _, nextEvent := range toEvents(forUser) {
		_, isAUserEditorEvent := nextEvent.(events.UserPromotedToEditor)
		if isAUserEditorEvent {
			return true, nil
		}
	}

	return false, nil
}
