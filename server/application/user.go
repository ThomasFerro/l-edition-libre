package application

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/contexts"
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

func IsAnEditor(ctx context.Context) (bool, error) {
	getHistory := contexts.FromContextOrDefault(ctx, contexts.ContextualizedUserHistoryContextKey{}, func(c context.Context) ([]ContextualizedEvent, error) {
		return []ContextualizedEvent{}, nil
	})
	history, err := getHistory(ctx)

	fmt.Printf("\n\n\n err %v \n history %v \n\n\n\n", err, history)
	if err != nil {
		return false, err
	}
	for _, nextEvent := range history {
		_, isAUserEditorEvent := nextEvent.Event().(events.UserPromotedToEditor)
		if isAUserEditorEvent {
			return true, nil
		}
	}

	return false, nil
}

func UserHaveAccessToManuscript(ctx context.Context) (bool, error) {
	isAnEditor, err := IsAnEditor(ctx)
	if err != nil {
		slog.Warn("user role check error", "error", err)
		return false, err
	}
	if isAnEditor {
		return true, nil
	}
	isManuscriptWriter, err := isTheManuscriptWriter(ctx)
	if err != nil {
		slog.Warn("user's link to manuscript check error", "error", err)
		return false, err
	}
	if isManuscriptWriter {
		return true, nil
	}
	slog.Warn("user is not the writer nor an editor")
	return false, nil
}

type UsersHistory interface {
	For(UserID) ([]ContextualizedEvent, error)
	Append(context.Context, []ContextualizedEvent) error
}
