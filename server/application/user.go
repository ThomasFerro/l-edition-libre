package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
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

func isAnEditor(history UsersHistory, userID UserID) (bool, error) {
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

func (app Application) UserHaveAccessToManuscript(userID UserID, manuscriptID ManuscriptID) (bool, error) {
	isAnEditor, err := isAnEditor(app.usersHistory, userID)
	if err != nil {
		slog.Warn("user role check error", "user_id", userID, "manuscript_id", manuscriptID, "error", err)
		return false, err
	}
	if isAnEditor {
		return true, nil
	}
	isManuscriptWriter, err := isTheManuscriptWriter(app.manuscriptsHistory, userID, manuscriptID)
	if err != nil {
		slog.Warn("user's link to manuscript check error", "user_id", userID, "manuscript_id", manuscriptID, "error", err)
		return false, err
	}
	if isManuscriptWriter {
		return true, nil
	}
	slog.Warn("user is not the writer nor an editor", "user_id", userID, "manuscript_id", manuscriptID)
	return false, nil
}

func (app Application) UserIsAnEditor(userID UserID) (bool, error) {
	return isAnEditor(app.usersHistory, userID)
}

type UsersHistory interface {
	For(UserID) ([]ContextualizedEvent, error)
	Append(context.Context, []ContextualizedEvent) error
}
